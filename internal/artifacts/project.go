package artifacts

import (
	"fmt"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/exec"
	"github.com/SAP/cloud-mta/mta"
)

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
	return execProjectBuilders(oMta, phase)
}

func execProjectBuilders(oMta *mta.MTA, phase string) error {
	if phase == "pre" && oMta.BuildParams != nil {
		return execProjectBuilder(oMta.BuildParams.BeforeAll)
	}
	if phase == "post" && oMta.BuildParams != nil {
		return execProjectBuilder(oMta.BuildParams.AfterAll)
	}
	return nil
}

func execProjectBuilder(builders []mta.ProjectBuilder) error {
	for _, builder := range builders {
		builderCommands, err := getProjectBuilderCommands(builder)
		if err != nil {
			return err
		}
		cmds := commands.CmdConverter(".", builderCommands.Command)
		// Execute commands
		err = exec.Execute(cmds)
		if err != nil {
			return err
		}
	}
	return nil
}

func getProjectBuilderCommands(builder mta.ProjectBuilder) (commands.CommandList, error) {
	dummyModule := mta.Module{}
	dummyModule.BuildParams = make(map[string]interface{})
	dummyModule.BuildParams["builder"] = builder.Builder
	dummyModule.BuildParams["commands"] = builder.Commands
	builderCommands, _, err := commands.CommandProvider(dummyModule)
	return builderCommands, err
}
