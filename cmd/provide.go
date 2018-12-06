package commands

import (
	"fmt"

	"cloud-mta-build-tool/internal/buildops"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"cloud-mta-build-tool/internal/fsys"
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
		err := dir.ValidateDeploymentDescriptor(descriptorPModuleFlag)
		if err == nil {
			ep := locationParameters(sourceBModuleFlag, targetBModuleFlag, descriptorPModuleFlag)
			err = provideModules(&ep)
		}
		if err != nil {
			err = errors.Wrap(err, "Modules provider failed")
		}
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func provideModules(file dir.IMtaParser) error {
	// read MTA from mta.yaml
	m, err := file.ParseFile()
	if err != nil {
		return err
	}
	modules, err := buildops.GetModulesNames(m)
	if err != nil {
		return err
	}
	// Get list of modules names
	fmt.Println(modules)
	return nil
}
