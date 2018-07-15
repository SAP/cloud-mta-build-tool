package commands

import (
	"github.com/spf13/cobra"

	"mbtv2/cmd/builders"
	fs "mbtv2/cmd/fsys"
	"mbtv2/cmd/proc"
)

// TODO - inject from outside
// Build for UI5 application
var html5 = &cobra.Command{
	Use:   "html5",
	Short: "Build for HTML5/UI5 project",
	Run: func(cmd *cobra.Command, args []string) {

		// Get MTA structure
		mtaStruct := proc.GetMta(fs.GetPath())
		// Read json configuration file
		cfg := proc.ReadConfig()
		// TODO - fetch the specific module
		for _, mod := range mtaStruct.Modules {
			switch mod.Type {
			case "html5":
				builders.Build(builders.NewGruntBuilder(mod.Path, mod.Name, cfg.TmpPath),
					fs.GetPath(), fs.DefaultTempDirFunc(fs.GetPath()))

			}
		}

	},
}
