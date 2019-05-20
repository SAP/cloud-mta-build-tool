package commands

import (
	"github.com/SAP/cloud-mta-build-tool/internal/exec"
	"os"

	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
	"github.com/SAP/cloud-mta-build-tool/internal/tpl"
)

const (
	makefile    = "Makefile.mta"
	makefileTmp = "Makefile_tmp.mta"
)

// flags of init command
var initCmdSrc string
var initCmdTrg string
var initCmdDesc string
var initCmdName string
var initCmdMode string

// flags of build command
var buildProjectCmdSrc string
var buildProjectCmdTrg string
var buildProjectCmdDesc string
var buildProjectCmdMode string
var buildProjectCmdMtar string
var buildProjectCmdPlatform string
var buildProjectCmdStrict bool

// init flags of init command
func init() {
	// set flags of init command
	initCmd.Flags().StringVarP(&initCmdSrc, "source", "s", "",
		"the path to the MTA project; the current path is set as the default")
	initCmd.Flags().StringVarP(&initCmdTrg, "target", "t", "",
		"the path to the MBT results folder; the current path is set as the default")
	initCmd.Flags().StringVarP(&initCmdDesc, "desc", "d", "",
		`the MTA descriptor; supported values: "dev" (development descriptor, default value) and "dep" (deployment descriptor)`)
	initCmd.Flags().StringVarP(&initCmdName, "name", "n", "",
		"the name of the Makefile; Makefile.mta is set as the default")
	initCmd.Flags().StringVarP(&initCmdMode, "mode", "m", "",
		"Mode of Makefile generation - default/verbose")

	// set flags of build command
	buildCmd.Flags().StringVarP(&buildProjectCmdSrc, "source", "s", "",
		"the path to the MTA project; the current path is set as the default")
	buildCmd.Flags().StringVarP(&buildProjectCmdTrg, "target", "t", "$(CURDIR)/mta_archives",
		"the path to the MBT results folder; the current path is set as the default")
	buildCmd.Flags().StringVarP(&buildProjectCmdDesc, "desc", "d", "",
		`the MTA descriptor; supported values: "dev" (development descriptor, default value) and "dep" (deployment descriptor)`)
	buildCmd.Flags().StringVarP(&buildProjectCmdMtar, "mtar", "", "*",
		"Mode of Makefile generation - default/verbose")
	buildCmd.Flags().StringVarP(&buildProjectCmdPlatform, "platform", "p", "",
		`the deployment platform; supported plaforms: "cf", "xsa", "neo"`)
	buildCmd.Flags().BoolVarP(&buildProjectCmdStrict, "strict", "", true,
		"true - duplicated fields and fields that are not defined reported as errors. false -  they are reported as warnings; true is set as a default")
}

// Generates the Makefile.mta file according to the MTA descriptor
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generates Makefile",
	Long:  "generates Makefile as a manifest file that describes the build process",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Generate build script
		err := tpl.ExecuteMake(initCmdSrc, initCmdTrg, makefile, initCmdDesc, initCmdMode, os.Getwd)
		logError(err)
	},
}

// Generates the Makefile.mta file according to the MTA descriptor and executes it
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "generates and executes Makefile",
	Long:  "generates Makefile as a manifest file that describes the build process and executes it",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Generate build script
		err := artifacts.ExecBuild(buildProjectCmdSrc, buildProjectCmdTrg, buildProjectCmdDesc, buildProjectCmdMode, buildProjectCmdMtar, buildProjectCmdPlatform, buildProjectCmdStrict, os.Getwd, exec.Execute)
		logError(err)
		return err
	},
}
