package commands

import (
	"fmt"

	"cloud-mta-build-tool/mta"
	"github.com/spf13/cobra"
)

// Provide list of modules
var pModuleCmd = &cobra.Command{
	Use:   "modules",
	Short: "Provide list of modules",
	Long:  "Provide list of modules",
	Run: func(cmd *cobra.Command, args []string) {
		// Get Yaml file
		err := provideModules("")
		LogError(err)
	},
}

func provideModules(path string) error {

	mo, err := mta.ReadMta(path, "mta.yaml")
	if err == nil {
		// Get list of modules names
		fmt.Println(mo.GetModulesNames())
	}
	return err
}
