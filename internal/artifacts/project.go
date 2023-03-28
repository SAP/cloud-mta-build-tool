package artifacts

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/exec"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"github.com/SAP/cloud-mta-build-tool/internal/tpl"
	"github.com/SAP/cloud-mta-build-tool/internal/version"
	"github.com/SAP/cloud-mta/mta"
)

const (
	copyInParallel = false
	// MaxMakeParallel - Maximum number of parallel makefile jobs if the parameter is not set by the user
	MaxMakeParallel = 8
)

// ExecBuild - Execute MTA project build
func ExecBuild(makefileTmp, source, target string, extensions []string, mode, mtar, platform string,
	strict bool, jobs int, outputSync bool, wdGetter func() (string, error), wdExec func([][]string, bool) error,
	useDefaultMbt bool, keepMakefile bool, sBomFilePath string) error {
	message, err := version.GetVersionMessage()
	if err == nil {
		logs.Logger.Info(message)
	}

	// Generate build script
	err = tpl.ExecuteMake(source, "", extensions, makefileTmp, mode, wdGetter, useDefaultMbt)
	if err != nil {
		return err
	}

	cmdParams := createMakeCommand(makefileTmp, source, target, mode, mtar, platform, strict, jobs,
		outputSync, runtime.NumCPU, sBomFilePath)
	err = wdExec([][]string{cmdParams}, false)

	// Remove temporary Makefile
	var removeError error = nil
	if !keepMakefile {
		removeError = os.Remove(filepath.Join(source, filepath.FromSlash(makefileTmp)))
		if removeError != nil {
			removeError = errors.Wrapf(removeError, removeFailedMsg, makefileTmp)
		}
	}

	if err != nil {
		if removeError != nil {
			logs.Logger.Error(removeError)
		}
		return errors.Wrap(err, execFailedMsg)
	}
	return removeError
}

func createMakeCommand(makefileName, source, target, mode, mtar, platform string, strict bool, jobs int,
	outputSync bool, numCPUGetter func() int, sBomFilePath string) []string {
	cmdParams := []string{source, "make", "-f", makefileName, "p=" + platform, "mtar=" + mtar, "strict=" + strconv.FormatBool(strict), "mode=" + mode}
	if target != "" {
		cmdParams = append(cmdParams, `t="`+target+`"`)
	}
	if tpl.IsVerboseMode(mode) {
		if jobs <= 0 {
			jobs = numCPUGetter()
			if jobs > MaxMakeParallel {
				jobs = MaxMakeParallel
			}
		}
		cmdParams = append(cmdParams, fmt.Sprintf("-j%d", jobs))

		if outputSync {
			cmdParams = append(cmdParams, "-Otarget")
		}
	}

	if strings.TrimSpace(source) != "" {
		cmdParams = append(cmdParams, `source="`+source+`"`)
	}

	if strings.TrimSpace(sBomFilePath) != "" {
		cmdParams = append(cmdParams, `sbom-file-path="`+sBomFilePath+`"`)
	}

	return cmdParams
}

// ExecuteProjectBuild - execute pre or post phase of project build
func ExecuteProjectBuild(source, target, descriptor string, extensions []string, phase string, getWd func() (string, error)) error {
	if phase != "pre" && phase != "post" {
		return fmt.Errorf(UnsupportedPhaseMsg, phase)
	}
	loc, err := dir.Location(source, target, descriptor, extensions, getWd)
	if err != nil {
		return err
	}
	oMta, err := loc.ParseFile()
	if err != nil {
		return err
	}
	return execProjectBuilders(loc, oMta, phase)
}

func execProjectBuilders(loc *dir.Loc, oMta *mta.MTA, phase string) error {
	if phase == "pre" && oMta.BuildParams != nil {
		return execProjectBuilder(oMta.BuildParams.BeforeAll, "before-all")
	}
	if phase == "post" {
		err := copyResourceContent(loc.GetSource(), loc.GetTargetTmpDir(), oMta, copyInParallel)
		if err != nil {
			return err
		}
		err = copyRequiredDependencyContent(loc.GetSource(), loc.GetTargetTmpDir(), oMta, copyInParallel)
		if err != nil {
			return err
		}
		if oMta.BuildParams != nil {
			return execProjectBuilder(oMta.BuildParams.AfterAll, "after-all")
		}
	}
	return nil
}

func execProjectBuilder(builders []mta.ProjectBuilder, phase string) error {
	errMessage := `the "%s"" build failed`
	logs.Logger.Infof(`running the "%s" build...`, phase)
	for _, builder := range builders {
		builderCommands, err := getProjectBuilderCommands(builder)
		if err != nil {
			return errors.Wrapf(err, errMessage, phase)
		}
		cmds, err := commands.CmdConverter(".", builderCommands.Command)
		if err != nil {
			return errors.Wrapf(err, errMessage, phase)
		}
		// Execute commands
		err = exec.ExecuteWithTimeout(cmds, builder.Timeout, true)
		if err != nil {
			return errors.Wrapf(err, errMessage, phase)
		}
	}
	return nil
}

func getProjectBuilderCommands(builder mta.ProjectBuilder) (commands.CommandList, error) {
	dummyModule := mta.Module{}
	dummyModule.BuildParams = make(map[string]interface{})
	dummyModule.BuildParams["builder"] = builder.Builder
	dummyModule.BuildParams["commands"] = builder.Commands
	if builder.Builder == "custom" && builder.Commands == nil && len(builder.Commands) == 0 {
		logs.Logger.Warn(commandsMissingMsg)
		return commands.CommandList{Command: []string{}}, nil
	}
	if builder.Builder != "custom" && builder.Commands != nil && len(builder.Commands) != 0 {
		logs.Logger.Warnf(commandsNotSupportedMsg, builder.Builder)
	}
	builderCommands, _, err := commands.CommandProvider(dummyModule)
	return builderCommands, err
}
