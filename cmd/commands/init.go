package commands

import (
	"cloud-mta-build-tool/cmd/gen"

	"github.com/spf13/cobra"
)

var initProcess = &cobra.Command{
	Use:   "init",
	Short: "Generate Makefile",
	Long:  "Generate Makefile as manifest which describe's the build process",
	Run: func(cmd *cobra.Command, args []string) {
		//Todo - remove the script option
		if (len(args) > 0) && (stringInSlice("script", args)) {
			// Generate build script
			if args[0] == "script" {
				gen.Generate("test")
			}
		} else {
			//Generate make
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
