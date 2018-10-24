package commands

import (
	"fmt"

	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/internal/version"
	"github.com/alecthomas/gometalinter/_linters/src/gopkg.in/yaml.v2"
	"github.com/spf13/cobra"
)

var buildTargetFlag string
var validationFlag string

func init() {

	// Add command to the root
	rootCmd.AddCommand(provideCmd, buildCmd, executeCmd, initProcessCmd, versionCmd)
	// Build module
	provideCmd.AddCommand(pModuleCmd)
	// Provide module
	buildCmd.AddCommand(bModuleCmd)
	// execute immutable commands
	executeCmd.AddCommand(packCmd, genMetaCmd, genMtadCmd, genMtarCmd, cleanupCmd, validateCmd)
	// build command target flags
	buildCmd.Flags().StringVarP(&buildTargetFlag, "target", "t", "", "Build for specified environment ")
	// validation flags , can be used for multiple scenario
	validateCmd.Flags().StringVarP(&validationFlag, "mode", "m", "", "Validation mode ")
}

// Parent command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build Project",
	Long:  "Build MTA project",
	Run:   nil,
}

// Parent command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute step",
	Long:  "Execute standalone step as part of the build process",
	Run:   nil,
}

// Parent command
var provideCmd = &cobra.Command{
	Use:   "provide",
	Short: "MBT data provider",
	Long:  "MBT data provider",
	Run:   nil,
}

// Parent command
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
