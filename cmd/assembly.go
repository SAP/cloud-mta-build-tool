package commands

import (
	"os"

	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
	"github.com/spf13/cobra"
)

const (
	defaultDeploymentDescriptor string = "dep"
	defaultPlatform             string = "cf"
	defaultMtaLocation          string = ""
	defaultMtaAssemblyLocation  string = ""
)

func init() {}

// Generate mtar from build artifacts
var assemblyCommand = &cobra.Command{
	Use:       "assemble",
	Short:     "Assemble MTA Archive",
	Long:      "Assemble MTA Archive",
	ValidArgs: []string{"Deployment descriptor location"},
	Args:      cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := artifacts.CopyMtaContent(defaultMtaLocation, defaultMtaAssemblyLocation, defaultDeploymentDescriptor, os.Getwd)
		if err != nil {
			logError(err)
			return err
		}
		err = artifacts.ExecuteGenMeta(defaultMtaLocation, defaultMtaAssemblyLocation, defaultDeploymentDescriptor, defaultPlatform, os.Getwd)
		if err != nil {
			logError(err)
			return err
		}
		err = artifacts.ExecuteGenMtar(defaultMtaLocation, defaultMtaAssemblyLocation, defaultDeploymentDescriptor, os.Getwd)
		if err != nil {
			logError(err)
			return err
		}

		err = artifacts.ExecuteCleanup(defaultMtaLocation, defaultMtaAssemblyLocation, defaultDeploymentDescriptor, os.Getwd)
		if err != nil {
			logError(err)
			return err
		}

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: false,
}
