package commands

import (
	"os"
	"time"

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
var buildProjectCmdMtar = "*"
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
	buildCmd.Flags().StringVarP(&buildProjectCmdTrg, "target", "t", "", "the path to the MBT results folder; the current path is set as the default")
	buildCmd.Flags().StringVarP(&buildProjectCmdMtar, "mtar", "", "", "the file name of the generated archive file")
	buildCmd.Flags().StringVarP(&buildProjectCmdPlatform, "platform", "p", "cf", `the deployment platform; supported platforms: "cf" (default value), "xsa", "neo"`)
	buildCmd.Flags().BoolVarP(&buildProjectCmdStrict, "strict", "", true, `if set to true, duplicated fields and fields not defined in the "mta.yaml" schema are reported as errors;
if set to false, they are reported as warnings`)
	buildCmd.Flags().BoolP("help", "h", false, `prints detailed information about the "build" command`)
	buildCmd.MarkFlagRequired("platform")
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
	Short: "Executes MTA project build",
	Long:  "Generates a temporary `Makefile` according to the MTA descriptor and runs the `make` command to package the MTA project into the MTA archive",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Generate temp Makefile with unique id
		makefileTmp := "Makefile_" + time.Now().Format("20060102150405") + ".mta"
		// Generate build script
		err := artifacts.ExecBuild(makefileTmp, buildProjectCmdSrc, buildProjectCmdTrg, "", buildProjectCmdMtar, buildProjectCmdPlatform, buildProjectCmdStrict, os.Getwd, exec.Execute)
		return err
	},
	SilenceUsage: true,
}
