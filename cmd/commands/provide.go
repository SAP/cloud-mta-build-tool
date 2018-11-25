package commands

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	fs "cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/mta"
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
		err := fs.ValidateDeploymentDescriptor(descriptorPModuleFlag)
		if err == nil {
			ep := locationParameters(sourceBModuleFlag, targetBModuleFlag, descriptorPModuleFlag)
			err = provideModules(&ep)
		}
		err = errors.Wrap(err, "Modules provider failed")
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func provideModules(ep *fs.Loc) error {
	// read MTA from mta.yaml
	m, err := mta.ParseFile(ep)
	if err != nil {
		return err
	}
	modules, err := m.GetModulesNames()
	if err != nil {
		return err
	}
	// Get list of modules names
	fmt.Println(modules)
	return nil
}
