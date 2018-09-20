package commands

import (
	"github.com/spf13/cobra"
)

// Build the whole MTA project as monolith
var cfBuild = &cobra.Command{
	Use:   "cf",
	Short: "Build to CF env",
	Long:  "Build to CF env",
	Run: func(cmd *cobra.Command, args []string) {

		// Todo support CF build

	},
}

var neoBuild = &cobra.Command{
	Use:   "neo",
	Short: "Build to Neo",
	Long:  "Build to Neo",
	Run: func(cmd *cobra.Command, args []string) {
		// Todo support neo build
	},
}
