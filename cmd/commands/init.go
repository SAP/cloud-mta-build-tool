package commands

import (
	"cloud-mta-build-tool/internal/fsys"
	"github.com/spf13/cobra"

	"cloud-mta-build-tool/internal/tpl"
)

var initModeFlag string
var descriptorInitFlag string
var sourceInitFlag string
var targetInitFlag string

func init() {
	initProcessCmd.Flags().StringVarP(&initModeFlag, "mode", "m", "", "Mode of Makefile generation - default/verbose")
	initProcessCmd.Flags().StringVarP(&descriptorInitFlag, "desc", "d", "", "Descriptor MTA - dev/dep")
	initProcessCmd.Flags().StringVarP(&sourceInitFlag, "source", "s", "", "Provide MTA source")
	initProcessCmd.Flags().StringVarP(&targetInitFlag, "target", "t", "", "Provide MTA target")
}

var initProcessCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate Makefile",
	Long:  "Generate Makefile as manifest which describe's the build process",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Generate build script
		err := dir.ValidateDeploymentDescriptor(descriptorInitFlag)
		if err == nil {
			err = tpl.Make(GetEndPoints(sourceInitFlag, targetInitFlag, descriptorInitFlag), initModeFlag)
		}
		LogErrorExt(err, "Makefile Generation failed")
	},
}
