package commands

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/mta/provider"
)

// Provide list of modules
var pm = &cobra.Command{
	Use:   "modules",
	Short: "Provide list of modules",
	Long:  "Provide list of modules",
	Run: func(cmd *cobra.Command, args []string) {
		names := provider.GetModulesNames(filepath.Join(fs.GetPath(), ""))
		fmt.Println(names)
	},
}
