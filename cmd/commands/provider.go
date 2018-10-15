package commands

import (
	"fmt"

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
		if err == nil {
			// Get list of modules names
			names := mo.GetModulesNames()
			if err == nil {
				fmt.Println(names)
			}
		}
		LogError(err)
	},
}
