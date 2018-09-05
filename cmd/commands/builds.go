package commands

import (
	"github.com/spf13/cobra"

	"cloud-mta-build-tool/cmd/build-executers"
	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/proc"
)

// TODO - Only for testing, will be removed
// Build for UI5 application
var html5 = &cobra.Command{
	Use:   "html5",
	Short: "Build for HTML5/UI5 project",
	Run: func(cmd *cobra.Command, args []string) {

		// Get MTA structure
		mtaStruct, err := proc.GetMta(fs.GetPath())
		if err != nil {
			//TODO error handling
		}
		// Read json configuration file
		cfg := proc.ReadConfig()
		// TODO - fetch the specific module
		for _, mod := range mtaStruct.Modules {
			switch mod.Type {
			case "html5":
				builders.Build(builders.NewGruntBuilder(mod.Path, mod.Name, cfg.TmpPath),
					fs.GetPath())

			}
		}

	},
}
