package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
	"github.com/SAP/cloud-mta-build-tool/internal/exec"
	"github.com/SAP/cloud-mta-build-tool/internal/tpl"
)

const (
	makefile = "Makefile.mta"
)

// flags of init command
var initCmdSrc string
var initCmdTrg string
var initCmdName string
var initCmdMode string

// flags of build command
var buildProjectCmdSrc string
var buildProjectCmdTrg string
var buildProjectCmdMtar string
var buildProjectCmdPlatform string
var buildProjectCmdStrict bool

func init() {
	// set flags for init command
	initCmd.Flags().StringVarP(&initCmdSrc, "source", "s", "", "the path to the MTA project; the current path is set as the default")
	initCmd.Flags().StringVarP(&initCmdTrg, "target", "t", "", "the path to the generated Makefile folder; the current path is set as the default")
	initCmd.Flags().StringVarP(&initCmdMode, "mode", "m", "", `the mode of the Makefile generation; supported values: "default" and "verbose"`)
	initCmd.Flags().MarkHidden("mode")
	initCmd.Flags().BoolP("help", "h", false, `prints detailed information about the "init" command`)

	// set flags of build command
	buildCmd.Flags().StringVarP(&buildProjectCmdSrc, "source", "s", "", "the path to the MTA project; the current path is set as the default")
	buildCmd.Flags().StringVarP(&buildProjectCmdTrg, "target", "t", "$(CURDIR)/mta_archives", "the path to the MBT results folder; the current path is set as the default")
	buildCmd.Flags().StringVarP(&buildProjectCmdMtar, "mtar", "", "*", "the file name of the generated archive file")
	buildCmd.Flags().StringVarP(&buildProjectCmdPlatform, "platform", "p", "", `the deployment platform; supported plaforms: "cf", "xsa", "neo"`)
	buildCmd.Flags().BoolVarP(&buildProjectCmdStrict, "strict", "", true, `if set to true, duplicated fields and fields not defined in the "mta.yaml" schema are reported as errors; if set to false, they are reported as warnings`)
	buildCmd.Flags().BoolP("help", "h", false, `prints detailed information about the "build" command`)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generates Makefile",
	Long:  "generates Makefile as a manifest file that describes the build process of the MTA project",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Generate build script
		err := tpl.ExecuteMake(initCmdSrc, initCmdTrg, makefile, initCmdMode, os.Getwd)
		logError(err)
	},
}

// Execute MTA project build
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Execute MTA project build",
	Long:  "Building each of the modules in the MTA project and archiving",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Generate build script
		err := artifacts.ExecBuild(buildProjectCmdSrc, buildProjectCmdTrg, "", buildProjectCmdMtar, buildProjectCmdPlatform, buildProjectCmdStrict, os.Getwd, exec.Execute)
		logError(err)
	},
}
