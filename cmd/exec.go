package commands

import (
	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/exec"
)

var executeCmdCommands []string
var executeCmdTimeout string
var executeCmdDir string
var copyCmdSrc string
var copyCmdTrg string
var copyCmdPatterns []string

// Execute commands in the current working directory with a timeout.
// This is used in verbose make files for implementing a timeout on module builds.
var executeCommand = &cobra.Command{
	Use:   "execute",
	Short: "Execute commands with timeout",
	Long:  "Execute commands with timeout",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := exec.ExecuteCommandsWithTimeout(executeCmdCommands, executeCmdTimeout, executeCmdDir, true)
		logError(err)
		return err
	},
	Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Copy files matching the specified patterns from the source path to the target path.
// This is used in verbose make files for copying artifacts from a module's dependencies before building the module.
var copyCmd = &cobra.Command{
	Use:   "cp",
	Short: "Copy files by patterns",
	Long:  "Copy files by patterns",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := dir.CopyByPatterns(copyCmdSrc, copyCmdTrg, copyCmdPatterns)
		logError(err)
		return err
	},
	Hidden:        true,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	// set flag of execute command
	executeCommand.Flags().StringArrayVarP(&executeCmdCommands,
		"commands", "c", nil, "Commands to run")
	executeCommand.Flags().StringVarP(&executeCmdTimeout,
		"timeout", "t", "", "The timeout after which the run stops, in the format [123h][123m][123s]; 10m is set as the default")
	executeCommand.Flags().StringVarP(&executeCmdDir,
		"dir", "d", "", "The path to the folder in which to execute the commands; the current path is set as the default")

	// set flags of copy command
	copyCmd.Flags().StringVarP(&copyCmdSrc, "source", "s", "",
		"The path to the source folder")
	copyCmd.Flags().StringVarP(&copyCmdTrg, "target", "t", "",
		"The path to the target folder")
	copyCmd.Flags().StringArrayVarP(&copyCmdPatterns,
		"patterns", "p", nil, "Patterns for matching the files and folders to copy")
}
