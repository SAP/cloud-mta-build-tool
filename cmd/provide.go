package commands

import (
	"os"

	"github.com/SAP/cloud-mta-build-tool/internal/buildops"

	"github.com/spf13/cobra"
)

var provideModuleCmdSrc string
var provideModuleCmdDesc string

// init - inits flags of provide module command
func init() {
	provideModuleCmd.Flags().StringVarP(&provideModuleCmdSrc, "source", "s", "", "Provide MTA source  ")
	provideModuleCmd.Flags().StringVarP(&provideModuleCmdDesc, "desc", "d", "", "Descriptor MTA - dev/dep")
}

// provideModuleCmd - Provide list of modules
var provideModuleCmd = &cobra.Command{
	Use:   "modules",
	Short: "Provide list of modules",
	Long:  "Provide list of modules",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := buildops.ProvideModules(provideModuleCmdSrc, provideModuleCmdDesc, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}
