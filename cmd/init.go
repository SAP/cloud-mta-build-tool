package commands

import (
	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
	"github.com/SAP/cloud-mta-build-tool/internal/tpl"
	"github.com/spf13/cobra"
	"os"
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
var buildProjectCmdName string
var buildProjectCmdMode string
var buildProjectCmdMtar string
var buildProjectCmdPlatform string
var buildProjectCmdStrict bool

// init flags of init command
func init() {
	// set flags of init command
	initCmd.Flags().StringVarP(&initCmdSrc, "source", "s", "", "Provide MTA source")
	initCmd.Flags().StringVarP(&initCmdTrg, "target", "t", "", "Provide MTA target")
	initCmd.Flags().StringVarP(&initCmdDesc, "desc", "d", "", "Descriptor MTA - dev/dep")
	initCmd.Flags().StringVarP(&initCmdDesc, "name", "n", "", "Name of Makefile")
	initCmd.Flags().StringVarP(&initCmdMode, "mode", "m", "", "Mode of Makefile generation - default/verbose")

	// set flags of build command
	buildCmd.Flags().StringVarP(&buildProjectCmdSrc, "source", "s", "$(CURDIR)", "Provide MTA source")
	buildCmd.Flags().StringVarP(&buildProjectCmdTrg, "target", "t", "$(CURDIR)/mta_archives", "Provide MTA target")
	buildCmd.Flags().StringVarP(&buildProjectCmdDesc, "desc", "d", "", "Descriptor MTA - dev/dep")
	buildCmd.Flags().StringVarP(&buildProjectCmdName, "name", "n", "", "Name of Makefile")
	buildCmd.Flags().StringVarP(&buildProjectCmdMode, "mode", "m", "", "Name of Mtar")
	buildCmd.Flags().StringVarP(&buildProjectCmdMtar, "mtar", "", "*", "Mode of Makefile generation - default/verbose")
	buildCmd.Flags().StringVarP(&buildProjectCmdPlatform, "platform", "p", "", "Platform")
	buildCmd.Flags().BoolVarP(&buildProjectCmdStrict, "strict", "", true, "Strict")
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
		err := artifacts.ExecBuild(buildProjectCmdSrc, buildProjectCmdTrg, buildProjectCmdDesc, buildProjectCmdMode, buildProjectCmdMtar, buildProjectCmdPlatform, buildProjectCmdStrict, os.Getwd)
		logError(err)
		return err
	},
}
