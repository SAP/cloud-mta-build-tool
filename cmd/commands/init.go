package commands

import (
	"github.com/spf13/cobra"
	"cloud-mta-build-tool/cmd/gen"
)

var initProcess = &cobra.Command{
	Use:   "init",
	Short: "Provide Build Script",
	Long:  "Provide Build Script",
	Run: func(cmd *cobra.Command, args []string) {
		//Todo - remove the script option
		if (len(args) > 0) && (stringInSlice("script", args)) {
			// Generate build script
			if args[0] == "script" {
				gen.Generate("test")
			}
		} else {
			gen.Make()
		}
	},
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
