package artifacts

import (
	"fmt"
	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/exec"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"github.com/SAP/cloud-mta-build-tool/internal/tpl"
	"github.com/SAP/cloud-mta/mta"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strconv"
)

const (
	copyInParallel = false
	makefileTmp    = "Makefile_tmp.mta"
)

// ExecBuild - Generates Makefile according to the MTA descriptor and executes it
func ExecBuild(buildProjectCmdSrc, buildProjectCmdTrg, buildProjectCmdDesc, buildProjectCmdMode, buildProjectCmdMtar, buildProjectCmdPlatform string, buildProjectCmdStrict bool, wdGetter func() (string, error)) error {
	// Generate build script
	err := tpl.ExecuteMake(buildProjectCmdSrc, "", makefileTmp, buildProjectCmdDesc, buildProjectCmdMode, wdGetter)
	if err != nil {
		return fmt.Errorf(`generation of the "%v" file failed`, makefileTmp)
	}
	err = exec.Execute([][]string{{buildProjectCmdSrc, "make", "-f", makefileTmp, " p=" + buildProjectCmdPlatform, " mtar=" + buildProjectCmdMtar, ` t="` + buildProjectCmdTrg + `"`, " strict=" + strconv.FormatBool(buildProjectCmdStrict)}})
	if err != nil {
		// Remove Makefile_tmp.mta file from directory
		err = os.Remove(filepath.FromSlash(makefileTmp))
		if err != nil {
			return fmt.Errorf(`removing of the "%v" file failed`, makefileTmp)
		}
		return fmt.Errorf(`execution of the "%v" file failed`, makefileTmp)
	}
	// Remove Makefile_tmp.mta file from directory
	err = os.Remove(filepath.FromSlash(makefileTmp))
	if err != nil {
		return fmt.Errorf(`removing of the "%v" file failed`, makefileTmp)
	}
	return err
}

// ExecuteProjectBuild - execute pre or post phase of project build
func ExecuteProjectBuild(source, descriptor, phase string, getWd func() (string, error)) error {
	if phase != "pre" && phase != "post" {
		return fmt.Errorf(`the "%s" phase of mta project build is invalid; supported phases: "pre", "post"`, phase)
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
		logs.Logger.Warn(`the "commands" property is missing in the "custom" builder`)
		return commands.CommandList{Command: []string{}}, nil
	}
	if builder.Builder != "custom" && builder.Commands != nil && len(builder.Commands) != 0 {
		logs.Logger.Warnf(`the "commands" property is not supported by the "%s" builder`, builder.Builder)
	}
	builderCommands, _, err := commands.CommandProvider(dummyModule)
	return builderCommands, err
}
