package commands

import (
	"fmt"

	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/mta"
	"github.com/spf13/cobra"
)

// Provide list of modules
var pModule = &cobra.Command{
	Use:   "modules",
	Short: "Provide list of modules",
	Long:  "Provide list of modules",
	Run: func(cmd *cobra.Command, args []string) {
		// Get Yaml file
		mo, err := mta.ReadMta("", "mta.yaml")
		if err != nil {
			logs.Logger.Error(err)
			return
		}
		// Get list of modules names
		names := mo.GetModulesNames()
		if err != nil {
			logs.Logger.Error(err)
			return
		}
		fmt.Println(names)
	},
}
