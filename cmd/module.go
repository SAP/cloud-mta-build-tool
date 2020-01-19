package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
)

// flags of pack command
var packCmdSrc string
var packCmdTrg string
var packCmdExtensions []string
var packCmdModule string
var packCmdPlatform string

// flags of build command
var buildModuleCmdSrc string
var buildModuleCmdTrg string
var buildModuleCmdExtensions []string
var buildModuleCmdModule string
var buildModuleCmdPlatform string

func init() {

	// sets the flags of of the command pack module
	packModuleCmd.Flags().StringVarP(&packCmdSrc, "source", "s", "",
		"The path to the MTA project; the current path is is set as default")
	packModuleCmd.Flags().StringVarP(&packCmdTrg, "target", "t", "",
		"The path to the folder in which the temporary artifacts of the module pack are created; the current path is is set as default")
	packModuleCmd.Flags().StringSliceVarP(&packCmdExtensions, "extensions", "e", nil,
		"The MTA extension descriptors")
	packModuleCmd.Flags().StringVarP(&packCmdModule, "module", "m", "",
		"The name of the module")
	packModuleCmd.Flags().StringVarP(&packCmdPlatform, "platform", "p", "cf",
		`The deployment platform; supported platforms: "cf", "xsa", "neo"`)

	// sets the flags of the command build module
	buildModuleCmd.Flags().StringVarP(&buildModuleCmdSrc, "source", "s", "",
		"The path to the MTA project; the current path is set as default")
	buildModuleCmd.Flags().StringVarP(&buildModuleCmdTrg, "target", "t", "",
		"The path to the folder in which the temporary artifacts of the module build are created; the current path is set as default")
	buildModuleCmd.Flags().StringSliceVarP(&buildModuleCmdExtensions, "extensions", "e", nil,
		"The MTA extension descriptors")
	buildModuleCmd.Flags().StringVarP(&buildModuleCmdModule, "module", "m", "",
		"The name of the module")
	buildModuleCmd.Flags().StringVarP(&buildModuleCmdPlatform, "platform", "p", "cf",
		`The deployment platform; supported platforms: "cf", "xsa", "neo"`)
}

// buildModuleCmd - Build module
var buildModuleCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds module",
	Long:  "Builds module and archives its artifacts",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteBuild(buildModuleCmdSrc, buildModuleCmdTrg, buildModuleCmdExtensions, buildModuleCmdModule, buildModuleCmdPlatform, os.Getwd)
		logError(err)
		return err
	},
	Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// zips the specific module and puts the artifacts in the temp folder according
// to the MTAR structure; that is, each module has new entry as folder in the MTAR folder
// Note - even if the path of the module was changed in the "mta.yaml" file, in the MTAR folder the
// the module folder gets the module name
var packModuleCmd = &cobra.Command{
	Use:   "pack",
	Short: "Packs module artifacts",
	Long:  "Packs the module artifacts after the build process",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecutePack(packCmdSrc, packCmdTrg, packCmdExtensions, packCmdModule, packCmdPlatform, os.Getwd)
		logError(err)
		return err
	},
	Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
}
