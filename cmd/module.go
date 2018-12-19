package commands

import (
	"os"

	"github.com/spf13/cobra"

	"cloud-mta-build-tool/internal/artifacts"
)

// flags of pack command
var packCmdSrc string
var packCmdTrg string
var packCmdDesc string
var packCmdModule string
var packCmdPlatform string

// flags of build command
var buildCmdSrc string
var buildCmdTrg string
var buildCmdDesc string
var buildCmdModule string
var buildCmdPlatform string

func init() {

	// set flags of command pack Module
	packModuleCmd.Flags().StringVarP(&packCmdSrc, "source", "s", "", "Provide MTA source ")
	packModuleCmd.Flags().StringVarP(&packCmdTrg, "target", "t", "", "Provide MTA target ")
	packModuleCmd.Flags().StringVarP(&packCmdDesc, "desc", "d", "", "Descriptor MTA - dev/dep")
	packModuleCmd.Flags().StringVarP(&packCmdModule, "module", "m", "", "Provide Module name ")
	packModuleCmd.Flags().StringVarP(&packCmdPlatform, "platform", "p", "", "Provide MTA platform ")

	// set flags of command build Module
	buildModuleCmd.Flags().StringVarP(&buildCmdSrc, "source", "s", "", "Provide MTA source ")
	buildModuleCmd.Flags().StringVarP(&buildCmdTrg, "target", "t", "", "Provide MTA target ")
	buildModuleCmd.Flags().StringVarP(&buildCmdDesc, "desc", "d", "", "Descriptor MTA - dev/dep")
	buildModuleCmd.Flags().StringVarP(&buildCmdModule, "module", "m", "", "Provide Module name ")
	buildModuleCmd.Flags().StringVarP(&buildCmdPlatform, "platform", "p", "", "Provide MTA platform ")
}

// buildModuleCmd - Build module
var buildModuleCmd = &cobra.Command{
	Use:   "build",
	Short: "Build module",
	Long:  "Build specific module according to the module name",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteBuild(buildCmdSrc, buildCmdTrg, buildCmdDesc, buildCmdModule, buildCmdPlatform, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: false,
}

// zip specific module and put the artifacts on the temp folder according
// to the mtar structure, i.e each module has new entry as folder in the mtar folder
// Note - even if the path of the module was changed in the mta.yaml in the mtar the
// the module folder will get the module name
var packModuleCmd = &cobra.Command{
	Use:   "pack",
	Short: "pack module artifacts",
	Long:  "pack the module artifacts after the build process",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecutePack(packCmdSrc, packCmdTrg, packCmdDesc, packCmdModule, packCmdPlatform, os.Getwd)
		logError(err)
		return err
	},
}
