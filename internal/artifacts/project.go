package artifacts

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/exec"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"github.com/SAP/cloud-mta-build-tool/internal/tpl"
	"github.com/SAP/cloud-mta/mta"
)

const (
	copyInParallel = false
)

// ExecBuild - Execute MTA project build
func ExecBuild(makefileTmp, buildProjectCmdSrc, buildProjectCmdTrg, buildProjectCmdMode, buildProjectCmdMtar, buildProjectCmdPlatform string, buildProjectCmdStrict bool, wdGetter func() (string, error), wdExec func([][]string) error) error {
	// Generate build script
	err := tpl.ExecuteMake(buildProjectCmdSrc, "", makefileTmp, buildProjectCmdMode, wdGetter)
	if err != nil {
		return errors.Wrapf(err, `generation of the "%s" file failed`, makefileTmp)
	}
	if buildProjectCmdTrg == "" {
		err = wdExec([][]string{{buildProjectCmdSrc, "make", "-f", makefileTmp, " p=" + buildProjectCmdPlatform, " mtar=" + buildProjectCmdMtar, " strict=" + strconv.FormatBool(buildProjectCmdStrict), " mode=" + buildProjectCmdMode}})
	} else {
		err = wdExec([][]string{{buildProjectCmdSrc, "make", "-f", makefileTmp, " p=" + buildProjectCmdPlatform, " mtar=" + buildProjectCmdMtar, ` t="` + buildProjectCmdTrg + `"`, " strict=" + strconv.FormatBool(buildProjectCmdStrict), " mode=" + buildProjectCmdMode}})
	}
	error := os.Remove(filepath.Join(buildProjectCmdSrc, filepath.FromSlash(makefileTmp)))
	if err != nil && error != nil {
		// Remove Makefile_tmp.mta file from directory
		return errors.Wrapf(err, execAndRemoveFailedMsg, makefileTmp, makefileTmp)
	} else if err != nil {
		return fmt.Errorf(execFailedMsg, makefileTmp)
	} else if error != nil {
		return fmt.Errorf(removeFailedMsg, makefileTmp)
	}
	return err
}

// ExecuteProjectBuild - execute pre or post phase of project build
func ExecuteProjectBuild(source, descriptor, phase string, getWd func() (string, error)) error {
	if phase != "pre" && phase != "post" {
		return fmt.Errorf(UnsupportedPhaseMsg, phase)
	}
	loc, err := dir.Location(source, "", descriptor, getWd)
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
		cmds := commands.CmdConverter(".", builderCommands.Command)
		// Execute commands
		err = exec.Execute(cmds)
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
