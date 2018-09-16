package commands

import (
	"github.com/spf13/cobra"

	"cloud-mta-build-tool/cmd/build-executers"
	"cloud-mta-build-tool/cmd/logs"
)

// Build the whole MTA project as monolith
var cfBuild = &cobra.Command{
	Use:   "cf",
	Short: "Build to CF env",
	Long:  "Build to CF env",
	Run: func(cmd *cobra.Command, args []string) {
		target := func(bld *builders.BuildCfg) {
			bld.Target = "cf"
		}
		_, err := builders.BuildProcess(target)
		if err != nil {
			logs.Logger.Error(err)
		}
	},
}

var neoBuild = &cobra.Command{
	Use:   "neo",
	Short: "Build to Neo",
	Long:  "Build to Neo",
	Run: func(cmd *cobra.Command, args []string) {
		// Todo support new build
	},
}
