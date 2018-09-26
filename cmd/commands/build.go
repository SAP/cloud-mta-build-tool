package commands

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"cloud-mta-build-tool/cmd/builders"
	"cloud-mta-build-tool/cmd/exec"
	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/mta"
	"cloud-mta-build-tool/mta/provider"
)

var buildTarget string

// Build module
var bm = &cobra.Command{
	Use:   "module",
	Short: "Build module",
	Long:  "Build module",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			buildModule(args[0])
		} else {
			logs.Logger.Errorf("Build module command require module name")
		}
	},
}

func buildModule(module string) {

	logs.Logger.Info("Start building module: ", module)
	mta, err := provider.MTA(filepath.Join(fs.GetPath(), ""))
	if err != nil {
		logs.Logger.Errorf("Error occurred while parsing the MTA file", err)
	}
	// Get module respective command's to execute
	mPathProp, mCmd := moduleCmd(mta, module)
	mRelPath := filepath.Join(fs.GetPath(), mPathProp)
	// Get module commands and path
	commands := cmdConverter(mRelPath, mCmd)
	// Execute child-process
	err = exec.Execute(commands)
	if err != nil {
		logs.Logger.Error(err)
	}

}

func init() {
	bm.Flags().StringVarP(&buildTarget, "target", "t", "", "Build for specified environment ")
}

// Get commands for specific module type
func moduleCmd(mta mta.MTA, moduleName string) (string, []string) {
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
