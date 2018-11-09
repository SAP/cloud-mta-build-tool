package commands

import (
	"fmt"

	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/internal/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {

	// Add command to the root
	rootCmd.AddCommand(provideCmd, executeCmd, initProcessCmd, versionCmd)
	// Build module
	provideCmd.AddCommand(pModuleCmd)
	// execute immutable commands
	executeCmd.AddCommand(bModuleCmd, packCmd, genMetaCmd, genMtadCmd, genMtarCmd, cleanupCmd, validateCmd)
	// build command target flags
}

// Parent command - Parent of all execution commands
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute step",
	Long:  "Execute standalone step as part of the build process",
	Run:   nil,
}

// Parent command - MTA info provider
var provideCmd = &cobra.Command{
	Use:   "provide",
	Short: "MBT data provider",
	Long:  "MBT data provider",
	Run:   nil,
}

// Parent command - CLI Version provider
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "MBT version",
	Long:  "MBT version",
	Run: func(cmd *cobra.Command, args []string) {
		err := printCliVersion()
		logError(err)
	},
}

func printCliVersion() error {
	v, err := version.GetVersion()
	if err == nil {
		fmt.Println(v.CliVersion)
	}
	return err
}

// logError - log errors if any
func logError(err error) {
	if err != nil {
		logs.Logger.Error(err)
	}
}

// logErrorExt - log error wrapped with new message
func logErrorExt(err error, newMsg string) {
	if err != nil {
		logs.Logger.Error(errors.Wrap(err, newMsg))
	}
}
