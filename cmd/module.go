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

// flags of Makefile build command
var buildModuleCmdSrc string
var buildModuleCmdTrg string
var buildModuleCmdExtensions []string
var buildModuleCmdModule string
var buildModuleCmdPlatform string

// flags of stand alone build command
var soloBuildModuleCmdSrc string
var soloBuildModuleCmdTrg string
var soloBuildModuleCmdExtensions []string
var soloBuildModuleCmdModules []string
var soloBuildModuleCmdAllDependencies bool

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

	// sets the flags of the Makefile command build module
	buildModuleCmd.Flags().StringVarP(&buildModuleCmdSrc, "source", "s", "",
		"The path to the MTA project; the current path is set as default")
	buildModuleCmd.Flags().StringVarP(&buildModuleCmdTrg, "target", "t", "",
		"The path to the folder in which the module build results are created; the <source folder>/.<project name>_mta_build_tmp/<module name> path is set as default")
	buildModuleCmd.Flags().StringSliceVarP(&buildModuleCmdExtensions, "extensions", "e", nil,
		"The MTA extension descriptors")
	buildModuleCmd.Flags().StringVarP(&buildModuleCmdModule, "module", "m", "",
		"The name of the module")
	buildModuleCmd.Flags().StringVarP(&buildModuleCmdPlatform, "platform", "p", "cf",
		`The deployment platform; supported platforms: "cf", "xsa", "neo"`)

	// sets the flags of the solo Makefile command build module
	soloBuildModuleCmd.Flags().StringVarP(&soloBuildModuleCmdSrc, "source", "s", "",
		"The path to the MTA project; the current path is set as default")
	soloBuildModuleCmd.Flags().StringVarP(&soloBuildModuleCmdTrg, "target", "t", "",
		"The path to the folder in which the module build results are created; the <current folder>/.<project name>_mta_build_tmp/<module name> path is set as default")
	soloBuildModuleCmd.Flags().StringSliceVarP(&soloBuildModuleCmdExtensions, "extensions", "e", nil,
		"The MTA extension descriptors")
	soloBuildModuleCmd.Flags().StringSliceVarP(&soloBuildModuleCmdModules, "modules", "m", nil,
		"The names of the modules")
	soloBuildModuleCmd.Flags().BoolVarP(&soloBuildModuleCmdAllDependencies, "with-all-dependencies", "a", false,
		"Build selected modules with all dependencies")
}

// soloBuildModuleCmd - Build module command used stand alone
var soloBuildModuleCmd = &cobra.Command{
	Use:   "module-build",
	Short: "Builds module according to configurations in the MTA development descriptor (mta.yaml)",
	Long:  "Builds module according to configurations in the MTA development descriptor (mta.yaml)",
	Args:  cobra.MaximumNArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteSoloBuild(soloBuildModuleCmdSrc, soloBuildModuleCmdTrg, soloBuildModuleCmdExtensions, soloBuildModuleCmdModules, soloBuildModuleCmdAllDependencies, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// buildModuleCmd - Build module command that is used in Makefile
var buildModuleCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds module",
	Long:  "Builds module according to configurations in the MTA development descriptor (mta.yaml) and archives its artifacts",
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
