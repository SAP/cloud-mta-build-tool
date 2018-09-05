package commands

import (
	"github.com/spf13/cobra"

	"cloud-mta-build-tool/cmd/gen"
	"cloud-mta-build-tool/cmd/logs"
)

var initProcess = &cobra.Command{
	Use:   "init",
	Short: "Generate Makefile",
	Long:  "Generate Makefile as manifest which describe's the build process",
	Run: func(cmd *cobra.Command, args []string) {
		// Todo - remove the script option
		if (len(args) > 0) && (stringInSlice("script", args)) {
			// Generate build script
			if args[0] == "script" {
				err := gen.Generate("test")
				if err != nil {
					logs.Logger.Error(err)
				}
			}
		} else {
			// Generate make
			if err := gen.Make(); err != nil {
				logs.Logger.Error(err)
			}

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
