package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

var cfgFile string

func init() {
	logs.NewLogger()
	cobra.OnInitialize(initConfig)
}

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "MBT",
	Short: "MTA Build tool",
	Long:  "MTA Build tool V2",
	Args:  cobra.MaximumNArgs(1),
}

// Execute command adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	logError(err)
}

func initConfig() {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logs.Logger.Println("Using config file:", viper.ConfigFileUsed())
	}
}
