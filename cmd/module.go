package commands

import (
	"github.com/spf13/cobra"
	"os"

	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
)

// flags of pack command
var packCmdSrc string
var packCmdTrg string
var packCmdDesc string
var packCmdModule string
var packCmdPlatform string

// flags of zip command
var zipCmdSrc string
var zipCmdTrg string
var zipCmdModule string
var zipCmdPlatform string

// flags of build command
var buildCmdSrc string
var buildCmdTrg string
var buildCmdDesc string
var buildCmdModule string
var buildCmdPlatform string

func init() {

	// sets the flags of of the command pack module
	packModuleCmd.Flags().StringVarP(&packCmdSrc, "source", "s", "",
		"the path to the MTA project; the current path is default")
	packModuleCmd.Flags().StringVarP(&packCmdTrg, "target", "t", "",
		"the path to the MBT results folder; the current path is default")
	packModuleCmd.Flags().StringVarP(&packCmdDesc, "desc", "d", "",
		"the MTA descriptor; supported values: dev (development descriptor, default value) and dep (deployment descriptor)")
	packModuleCmd.Flags().StringVarP(&packCmdModule, "module", "m", "",
		"the name of the module")
	packModuleCmd.Flags().StringVarP(&packCmdPlatform, "platform", "p", "",
		"the deployment platform; supported plaforms: cf, xsa, neo")

	// sets the flags of of the command zip module
	zipModuleCmd.Flags().StringVarP(&zipCmdSrc, "source", "s",
		"", "the path to the MTA project; the current path is default")
	zipModuleCmd.Flags().StringVarP(&zipCmdTrg, "target", "t", "",
		"the path to the MBT results folder; the current path is default")
	zipModuleCmd.Flags().StringVarP(&zipCmdModule, "module", "m", "",
		"the name of the module")

	// sets the flags of the command build module
	buildModuleCmd.Flags().StringVarP(&buildCmdSrc, "source", "s", "",
		"the path to the MTA project; the current path is default")
	buildModuleCmd.Flags().StringVarP(&buildCmdTrg, "target", "t", "",
		"the path to the MBT results folder; the current path is default")
	buildModuleCmd.Flags().StringVarP(&buildCmdDesc, "desc", "d", "",
		"the MTA descriptor; supported values: dev (development descriptor, default value) and dep (deployment descriptor)")
	buildModuleCmd.Flags().StringVarP(&buildCmdModule, "module", "m", "",
		"the name of the module")
	buildModuleCmd.Flags().StringVarP(&buildCmdPlatform, "platform", "p", "",
		"the deployment platform; supported plaforms: cf, xsa, neo")
}

// buildModuleCmd - Build module
var buildModuleCmd = &cobra.Command{
	Use:   "build",
	Short: "builds module",
	Long:  "builds module and archives its artifacts",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteBuild(buildCmdSrc, buildCmdTrg, buildCmdDesc, buildCmdModule, buildCmdPlatform, os.Getwd)
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
	Short: "packs module artifacts",
	Long:  "packs the module artifacts after the build process",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecutePack(packCmdSrc, packCmdTrg, packCmdDesc, packCmdModule, packCmdPlatform, os.Getwd)
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
var zipModuleCmd = &cobra.Command{
	Use:   "zip",
	Short: "zip module artifacts",
	Long:  "zip the module artifacts before the build process",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := artifacts.ExecuteZip(zipCmdSrc, zipCmdTrg, zipCmdModule, os.Getwd)
		logError(err)
		return err
	},
	Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
}
