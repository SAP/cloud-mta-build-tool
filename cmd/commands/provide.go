package commands

import (
	"fmt"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/mta"
	"github.com/spf13/cobra"
)

func init() {
	pModuleCmd.Flags().StringVarP(&pSourceFlag, "source", "s", "", "Provide MTA source  ")
}

// Provide list of modules
var pModuleCmd = &cobra.Command{
	Use:   "modules",
	Short: "Provide list of modules",
	Long:  "Provide list of modules",
	Run: func(cmd *cobra.Command, args []string) {
		err := provideModules(GetEndPoints())
		LogError(err)
	},
}

func provideModules(ep dir.EndPoints) error {
	// read MTA from mta.yaml
	mo, err := mta.ReadMta(ep)
	if err == nil {
		// Get list of modules names
		fmt.Println(mo.GetModulesNames())
	}
	return err
}
