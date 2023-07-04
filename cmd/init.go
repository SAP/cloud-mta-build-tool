package commands

import (
	"fmt"
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
var initCmdExtensions []string
var initCmdMode string

// flags of build command
var mbtCmdCLI string
var buildCmdSrc string
var buildCmdTrg string
var buildCmdExtensions []string
var buildCmdMtar = "*"
var buildCmdPlatform string
var buildCmdStrict bool
var buildCmdMode string
var buildCmdJobs int
var buildCmdOutputSync bool
var buildCmdKeepMakefile bool

func init() {
	// set flags for init command
	initCmd.Flags().StringVarP(&initCmdSrc, "source", "s", "", "The path to the MTA project; the current path is set as default")
	initCmd.Flags().StringVarP(&initCmdTrg, "target", "t", "", "The path to the folder in which the Makefile is generated; the current path is set as default")
	initCmd.Flags().StringSliceVarP(&initCmdExtensions, "extensions", "e", nil, "The MTA extension descriptors")
	initCmd.Flags().StringVarP(&initCmdMode, "mode", "m", "", `The mode of the Makefile generation; supported values: "default" and "verbose"`)
	_ = initCmd.Flags().MarkHidden("mode")
	initCmd.Flags().BoolP("help", "h", false, `Displays detailed information about the "init" command`)

	// set flags of build command
	buildCmd.Flags().StringVarP(&buildCmdSrc, "source", "s", "", "The path to the MTA project; the current path is set as default")
	buildCmd.Flags().StringVarP(&buildCmdTrg, "target", "t", "", `The path to the folder in which the MTAR file is created; the path to the "mta_archives" subfolder of the current folder is set as default`)
	buildCmd.Flags().StringSliceVarP(&buildCmdExtensions, "extensions", "e", nil, "The MTA extension descriptors")
	buildCmd.Flags().StringVarP(&buildCmdMtar, "mtar", "", "", "The file name of the generated archive file")
	buildCmd.Flags().StringVarP(&buildCmdPlatform, "platform", "p", "cf", `The deployment platform; supported platforms: "cf", "xsa", "neo"`)
	buildCmd.Flags().BoolVarP(&buildCmdStrict, "strict", "", true, `If set to true, duplicated fields and fields not defined in the "mta.yaml" schema are reported as errors; if set to false, they are reported as warnings`)
	buildCmd.Flags().StringVarP(&buildCmdMode, "mode", "m", "", `(beta) If set to "verbose", Make can run build jobs simultaneously.`)
	buildCmd.Flags().IntVarP(&buildCmdJobs, "jobs", "j", 0, fmt.Sprintf(`(beta) The number of Make jobs to be executed simultaneously. The default value is the number of available CPUs (maximum %d). Used only in "verbose" mode.`, artifacts.MaxMakeParallel))
	buildCmd.Flags().BoolVarP(&buildCmdOutputSync, "output-sync", "o", false, `(beta) Groups the output of each Make job and prints it when the job is complete. Used only in "verbose" mode.`)
	buildCmd.Flags().BoolVarP(&buildCmdKeepMakefile, "keep-makefile", "k", false, `Don't remove the generated Makefile after the build ends.`)
	_ = buildCmd.Flags().MarkHidden("keep-makefile")
	buildCmd.Flags().BoolP("help", "h", false, `Displays detailed information about the "build" command`)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generates a GNU Make manifest file that describes the build process of the MTA project",
	Long:  "Generates a GNU Make manifest file that describes the build process of the MTA project",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Generate build script
		err := tpl.ExecuteMake(initCmdSrc, initCmdTrg, initCmdExtensions, makefile, initCmdMode, os.Getwd, true)
		logError(err)
	},
}

// Execute MTA project build
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds the project modules and generates an MTA archive according to the MTA development descriptor (mta.yaml)",
	Long:  "Builds the project modules and generates an MTA archive according to the MTA development descriptor (mta.yaml)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Generate temp Makefile with unique id
		makefileTmp := "Makefile_" + time.Now().Format("20060102150405") + ".mta"
		// Generate build script
		// We want to use the current mbt and not the default (globally installed) mbt from the path when running mbt build, to allow users to run mbt build without the need to set the path.
		// This also supports using multiple versions of the mbt.
		// However, in some environments we might want to always use the default mbt from the path. This can be set by using environment variable MBT_USE_DEFAULT.
		useDefaultMbt := os.Getenv("MBT_USE_DEFAULT") == "true"
		// Note: we can only use the non-default mbt (i.e. the current executable name) from inside the command itself because if this function runs from other places like tests it won't point to the MBT
		err := artifacts.ExecBuild(makefileTmp, buildCmdSrc, buildCmdTrg, buildCmdExtensions, buildCmdMode, buildCmdMtar, buildCmdPlatform, buildCmdStrict, buildCmdJobs, buildCmdOutputSync, os.Getwd, exec.Execute, useDefaultMbt, buildCmdKeepMakefile)
		return err
	},
	SilenceUsage: true,
}
