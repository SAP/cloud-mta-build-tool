package commands

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"cloud-mta-build-tool/cmd/builders"
	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/cmd/mta/models"
	"cloud-mta-build-tool/cmd/mta/provider"
)

// Build module
var bm = &cobra.Command{
	Use:   "module",
	Short: "Build module",
	Long:  "Build module",
	Run: func(cmd *cobra.Command, args []string) {
		mta, err := provider.MTA(filepath.Join(fs.GetPath(), ""))
		if err != nil {
			logs.Logger.Errorf("Not able to parse MTA ", err)
		}
		module := moduleCmd(mta, args[0])
		logs.Logger.Info(module)
	},
}

func moduleCmd(mta models.MTA, moduleName string) []string {
	var cmd []string
	for _, m := range mta.Modules {
		if m.Name == moduleName {
			commandProvider := builders.CommandProvider(*m)
			cmd = commandProvider.Command
			break
		}
	}
	return cmd
}
