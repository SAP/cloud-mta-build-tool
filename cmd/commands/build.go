package commands

import (
	"path/filepath"
	"strings"

	"cloud-mta-build-tool/cmd/exec"

	"github.com/spf13/cobra"

	"cloud-mta-build-tool/cmd/builders"
	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/mta/models"
	"cloud-mta-build-tool/mta/provider"
)

// Build module
var bm = &cobra.Command{
	Use:   "module",
	Short: "Build module",
	Long:  "Build module",
	Run: func(cmd *cobra.Command, args []string) {
		// Get MTA
		mta, err := provider.MTA(filepath.Join(fs.GetPath(), ""))
		if err != nil {
			logs.Logger.Errorf("Not able to parse MTA ", err)
		}
		// Get module respective command's to execute
		mPathProp, mCmd := moduleCmd(mta, args[0])
		logs.Logger.Info(mPathProp, mCmd)
		// Get running project path
		pPath := fs.ProjectPath()
		mRelPath := filepath.Join(pPath, mPathProp)
		// Get module commands and path
		commands := cmdConverter(mRelPath, mCmd)
		logs.Logger.Info(commands)
		// Execute child-process
		err = exec.Execute(commands)
		if err != nil {
			logs.Logger.Error(err)
		}
	},
}

// Get commands for specific module type
func moduleCmd(mta models.MTA, moduleName string) (string, []string) {
	var cmd []string
	var mPath string
	for _, m := range mta.Modules {
		if m.Name == moduleName {
			commandProvider := builders.CommandProvider(*m)
			cmd = commandProvider.Command
			mPath = m.Path
			break
		}
	}
	return mPath, cmd
}

// Path and commands to execute
func cmdConverter(mPath string, cmdList []string) [][]string {
	var cmd [][]string
	for i := 0; i < len(cmdList); i++ {
		cmd = append(cmd, append([]string{mPath}, strings.Split(cmdList[i], " ")...))
	}
	return cmd
}
