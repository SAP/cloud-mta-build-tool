package commands

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"cloud-mta-build-tool/internal/builders"
	"cloud-mta-build-tool/internal/exec"
	fs "cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"
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
		err := buildModule(args[0])
		if err != nil {
			logs.Logger.Error(err)
		}
	},
}

func buildModule(module string) error {

	logs.Logger.Info("Start building module: ", module)
	// Read MTA Yaml File
	mta, err := mta.ReadMta("", "mta.yaml")
	if err == nil {
		// Get module respective command's to execute
		moduleRelPath, mCmd, err := moduleCmd(*mta, module)
		if err != nil {
			return err
		}
		modulePath, err := fs.GetFullPath(moduleRelPath)
		if err == nil {
			// Get module commands
			commands := cmdConverter(modulePath, mCmd)
			// Get temp dir for packing the artifacts
			artifactsPath, err := fs.GetArtifactsPath(modulePath)
			if err == nil {
				// Execute child-process with module respective commands
				err = exec.Execute(commands)
				if err == nil {
					// Pack the modules build artifacts (include node modules)
					// into the artifactsPath dir as data zip
					err = packModule(artifactsPath, moduleRelPath, module)
				}
			}
		}
	}
	return err
}

// Get commands for specific module type
func moduleCmd(mta mta.MTA, moduleName string) (string, []string, error) {
	var cmd []string
	var mPath string
	for _, m := range mta.Modules {
		if m.Name == moduleName {
			commandProvider, err := builders.CommandProvider(*m)
			if err != nil {
				return "", nil, err
			}
			cmd = commandProvider.Command
			mPath = m.Path
			break
		}
	}
	return mPath, cmd, nil
}

// path and commands to execute
func cmdConverter(mPath string, cmdList []string) [][]string {
	var cmd [][]string
	for i := 0; i < len(cmdList); i++ {
		cmd = append(cmd, append([]string{mPath}, strings.Split(cmdList[i], " ")...))
	}
	return cmd
}
