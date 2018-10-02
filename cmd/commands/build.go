package commands

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"cloud-mta-build-tool/cmd/builders"
	"cloud-mta-build-tool/cmd/exec"
	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/mta"
)

var buildTarget string

func init() {
	// Add environment flag for build purpose
	bModule.Flags().StringVarP(&buildTarget, "target", "t", "", "Build for specified environment ")
}

// Build module
var bModule = &cobra.Command{
	Use:   "module",
	Short: "Build module",
	Long:  "Build specific module according to the module name",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("build module command require module name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		buildModule(args[0])
	},
}

func buildModule(module string) {

	logs.Logger.Info("Start building module: ", module)

	pPath := fs.ProjectPath()
	// Read MTA file
	yamlFile, err := ioutil.ReadFile(pPath + pathSep + "mta.yaml")
	if err != nil {
		logs.Logger.Error("Not able to read the mta file ")
	}
	// Parse MTA file
	mta, err := mta.Parse(yamlFile)
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
