package commands

import (
	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/exec"
)

var executeCmdCommands []string
var executeCmdTimeout string

// Execute commands in the current working directory with a timeout.
// This is used in verbose make files for implementing a timeout on module builds.
var executeCommand = &cobra.Command{
	Use:   "execute",
	Short: "Execute commands with timeout",
	Long:  "Execute commands in the current working directory with timeout",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := exec.ExecuteCommandsWithTimeout(executeCmdCommands, executeCmdTimeout, true)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	Hidden:        true,
	SilenceErrors: true,
}

func init() {
	// set flag of execute command
	executeCommand.Flags().StringArrayVarP(&executeCmdCommands,
		"commands", "c", nil, "commands to run")
	executeCommand.Flags().StringVarP(&executeCmdTimeout,
		"timeout", "t", "", "the timeout after which the run stops, in the format [123h][123m][123s]; 10m is set as the default")
}
