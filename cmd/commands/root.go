package commands

import (
	"fmt"
	"os"

	"cloud-mta-build-tool/cmd/logs"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "mbt",
	Short: "MTA Build tool",
	Long:  "MTA Build tool V2",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main().
// It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// TODO -initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// TODO Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			logs.Logger.Error(err)
		}
		// Search config in home directory with name ".mbt" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mbt")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
