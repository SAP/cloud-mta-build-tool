package artifacts

import (
	"fmt"

	"github.com/SAP/cloud-mta/mta"
	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/exec"
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
	if phase == "pre" && oMta.BuildParams != nil {
		return execBuilder(beforeExec(oMta.BuildParams))
	}
	if phase == "post" && oMta.BuildParams != nil {
		return execBuilder(afterExec(oMta.BuildParams))
	}
	return nil
}

// get build params for before-all section
func beforeExec(pb *mta.ProjectBuild) string {
	for _, v := range pb.BeforeAll.Builders {
		return v.Builder
	}
	return ""
}

// get build params for after-all section
func afterExec(pb *mta.ProjectBuild) string {
	for _, v := range pb.AfterAll.Builders {
		return v.Builder
	}
	return ""
}

func execBuilder(builder string) error {
	if builder == "" {
		return nil
	}
	dummyModule := mta.Module{
		BuildParams: map[string]interface{}{
			"builder": builder,
		},
	}
	builderCommands, _, err := commands.CommandProvider(dummyModule)
	if err != nil {
		return errors.Wrap(err, "failed to parse the builder types configuration")
	}
	cmds := commands.CmdConverter(".", builderCommands.Command)
	// Execute commands
	err = exec.Execute(cmds)
	if err != nil {
		return errors.Wrapf(err, `the "%v" builder failed when executing commands`, builder)
	}
	return err
}
