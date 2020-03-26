package cmd

import (
	"os"

	"github.com/hellofresh/klepto/pkg/config"
	"github.com/hellofresh/klepto/pkg/formatter"
	wErrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	globalConfig   *config.Spec
	configFile     string
	configFileName = ".klepto.toml"
	verbose        bool

	// RootCmd steals and anonymises databases
	RootCmd = &cobra.Command{
		Use:   "klepto",
		Short: "Steals and anonymises databases",
		Long: `Klepto by HelloFresh.
		Takes the structure and data from one (mysql) database (--from),
		anonymises the data according to the provided configuration file,
		and inserts that data into another mysql database (--to).
	
		Perfect for bringing your live data to staging!`,
		Example: "klepto steal -c .klepto.toml|yaml|json --from root:root@localhost:3306/fromDb --to root:root@localhost:3306/toDb",
	}
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to config file (default is ./.klepto)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Make the operation more talkative")

	RootCmd.AddCommand(NewStealCmd())
	RootCmd.AddCommand(NewVersionCmd())
	RootCmd.AddCommand(NewUpdateCmd())
	RootCmd.AddCommand(NewInitCmd())

	log.SetOutput(os.Stderr)
	log.SetFormatter(&formatter.CliFormatter{})
}

func initConfig(c *cobra.Command, args []string) error {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugf("Reading config from %s...", configFileName)

	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return wErrors.Wrap(err, "can't find current working directory")
		}

		viper.SetConfigName(".klepto")
		viper.AddConfigPath(cwd)
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig()
	if err != nil {
		return wErrors.Wrap(err, "could not read configurations")
	}

	err = viper.Unmarshal(&globalConfig)
	if err != nil {
		return wErrors.Wrap(err, "could not unmarshal config file")
	}

	return nil
}
