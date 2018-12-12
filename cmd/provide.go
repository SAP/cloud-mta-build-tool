package commands

import (
	"os"

	"cloud-mta-build-tool/internal/build-ops"

	"github.com/spf13/cobra"
)

var sourcePModuleFlag string
var descriptorPModuleFlag string

func init() {
	pModuleCmd.Flags().StringVarP(&sourcePModuleFlag, "source", "s", "", "Provide MTA source  ")
	pModuleCmd.Flags().StringVarP(&descriptorPModuleFlag, "desc", "d", "", "Descriptor MTA - dev/dep")
}

// Provide list of modules
var pModuleCmd = &cobra.Command{
	Use:   "modules",
	Short: "Provide list of modules",
	Long:  "Provide list of modules",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := build_ops.ProvideModules(sourcePModuleFlag, descriptorPModuleFlag, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}
