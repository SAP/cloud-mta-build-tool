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

// Provide mtad.yaml from mta.yaml
var pMtad = &cobra.Command{
	Use:   "mtad",
	Short: "Provide mtad",
	Long:  "Provide deployment descriptor (mtad.yaml) from development descriptor (mta.yaml)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		mtaStr, err := mta.ReadMta(pMtadSourceFlag, "mta.yaml")
		if err == nil {
			err = mta.GenMtad(*mtaStr, pMtadTargetFlag, func(mtaStr mta.MTA) {
				convertTypes(mtaStr)
			})
		}
		LogError(err)
	},
}
