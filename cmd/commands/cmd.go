package commands

import (
	"fmt"

	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/internal/version"
	"github.com/alecthomas/gometalinter/_linters/src/gopkg.in/yaml.v2"
	"github.com/spf13/cobra"
)

var buildTargetFlag string

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
		LogError(err)
	},
}

func printCliVersion() error {
	v := version.Version{}
	err := yaml.Unmarshal(version.VersionConfig, &v)
	if err == nil {
		fmt.Println(v.CliVersion)
	}
	return err
}

// LogError - log errors if any
func LogError(err error) {
	if err != nil {
		logs.Logger.Error(err)
	}
}
