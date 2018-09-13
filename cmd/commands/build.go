package commands

import (
	"github.com/spf13/cobra"
)

// Build module
var bm = &cobra.Command{
	Use:   "module",
	Short: "Build module",
	Long:  "Build module",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
