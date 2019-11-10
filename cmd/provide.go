package commands

import (
	"os"

	"github.com/SAP/cloud-mta-build-tool/internal/buildops"

	"github.com/spf13/cobra"
)

var provideModuleCmdSrc string
var provideModuleCmdDesc string
var provideModuleCmdExtensions []string

// init - inits flags of provide module command
func init() {
	provideModuleCmd.Flags().StringVarP(&provideModuleCmdSrc, "source", "s",
		"", "The path to the MTA project; the current path is set as the default")
	provideModuleCmd.Flags().StringVarP(&provideModuleCmdDesc, "desc", "d", "",
		`The MTA descriptor; supported values: "dev" (development descriptor, default value) and "dep" (deployment descriptor)`)
	provideModuleCmd.Flags().StringSliceVarP(&provideModuleCmdExtensions, "extensions", "e", nil,
		"The MTA extension descriptors")
}

// provideModuleCmd - Provide list of modules
var provideModuleCmd = &cobra.Command{
	Use:   "modules",
	Short: "Provides list of modules",
	Long:  "Provides list of MTA project modules sorted by their dependencies",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := buildops.ProvideModules(provideModuleCmdSrc, provideModuleCmdDesc, provideModuleCmdExtensions, os.Getwd)
		logError(err)
		return err
	},
	Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
}
