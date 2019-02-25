package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/x-cray/logrus-prefixed-formatter"

	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

var cfgFile string

func init() {
	logs.Logger = logs.NewLogger()
	formatter, ok := logs.Logger.Formatter.(*prefixed.TextFormatter)
	if ok {
		formatter.DisableColors = true
	}
	cobra.OnInitialize(initConfig)
}

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:     "MBT",
	Short:   "MTA Build tool",
	Long:    "MTA Build tool V2",
	Version: cliVersion(),
	Args:    cobra.MaximumNArgs(1),
}

// Execute command adds all child commands to the root command and sets flags appropriately
func Execute() error {
	return rootCmd.Execute()
}

func initConfig() {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv() // reads in the environment variables that match
	// if a config file is found, reads it in
	if err := viper.ReadInConfig(); err == nil {
		logs.Logger.Println("Using config file:", viper.ConfigFileUsed())
	}
}
