package commands

import (
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/mta"
	"fmt"
	"github.com/spf13/cobra"
)

// Provide list of modules
var pModule = &cobra.Command{
	Use:   "modules",
	Short: "Provide list of modules",
	Long:  "Provide list of modules",
	Run: func(cmd *cobra.Command, args []string) {
		// Get Yaml file
		bytes, err := getFile(pathSep)
		if err != nil {
			logs.Logger.Error(err)
		}
		names, err := mta.GetModulesNames(bytes)
		if err != nil {
			logs.Logger.Error(err)
		}
		fmt.Println(names)
	},
}

// Provide mtad.yaml from mta.yaml
var pMtad = &cobra.Command{
	Use:   "mtad",
	Short: "Provide mtad",
	Long:  "Provide deployment descriptor (mtad.yaml) from development descriptor (mta.yaml)",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
