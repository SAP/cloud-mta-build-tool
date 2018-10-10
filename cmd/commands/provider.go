package commands

import (
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/mta"
	"fmt"
	"github.com/spf13/cobra"
)

var pMtadSourceFlag string
var pMtadTargetFlag string

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

// Provide mtad.yaml from mta.yaml
var pMtad = &cobra.Command{
	Use:   "mtad",
	Short: "Provide mtad",
	Long:  "Provide deployment descriptor (mtad.yaml) from development descriptor (mta.yaml)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		mtaStr, err := mta.ReadMta(pMtadSourceFlag, "mta.yaml")
		if err == nil {
			err = mta.GenMtad(*mtaStr, pMtadTargetFlag)
		}
		if err != nil {
			logs.Logger.Error(err)
		}
	},
}

func init() {
	pMtad.Flags().StringVarP(&pMtadSourceFlag, "source", "s", "", "Provide MTAD source ")
	pMtad.Flags().StringVarP(&pMtadTargetFlag, "target", "t", "", "Provide MTAD target ")
}
