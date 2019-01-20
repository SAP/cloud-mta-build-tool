package commands

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/buildops"
	"cloud-mta-build-tool/internal/fs"

	"github.com/SAP/cloud-mta/mta"
)

// CommandList - list of command to execute
type CommandList struct {
	Info    string
	Command []string
}

// CommandProvider - Get build command's to execute
//noinspection GoExportedFuncWithUnexportedType
func CommandProvider(modules mta.Module) (CommandList, error) {
	// Get config from ./commands_cfg.yaml as generated artifacts from source
	commands, err := parseModuleTypes(CommandsConfig)
	if err != nil {
		return CommandList{}, errors.Wrap(err, "failed to parse the commands configuration file")
	}
	customCommands, err := parseBuilders(CustomCommandsConfig)
	if err != nil {
		return CommandList{}, errors.Wrap(err, "failed to parse the custom commands configuration file")
	}
	return mesh(&modules, &commands, customCommands)
}

// Match the object according to type and provide the respective command
func mesh(module *mta.Module, moduleTypes *ModuleTypes, customCommands Builders) (CommandList, error) {
	// The object support deep struct for future use, can be simplified to flat object
	var cmds CommandList
	var commands []Commands
	var err error

	// get builder - module type name or custom builder if defined
	// and indicator if custom builder
	builder, custom := buildops.GetBuilder(module)

	// if module type used - get from module types configuration corresponding commands or custom builder if defined
	if !custom {
		for _, m := range moduleTypes.ModuleTypes {
			if m.Name == builder {
				if m.Builder != "" {
					// custom builder defined
					// check that no commands defined for module type
					if m.Type != nil && len(m.Type) > 0 {
						return cmds, fmt.Errorf("the module type definition can include builder or commands exclusively; the %s module type includes both", m.Name)
					}
					// continue with custom builders search
					builder = m.Builder
					custom = true
				} else {
					// get related information
					cmds.Info = m.Info
					commands = m.Type
				}
			}
		}
	}

	if custom {
		// custom builder used => get commands and info
		commands, cmds.Info, err = getCustomCommandsByBuilder(customCommands, builder)
		if err != nil {
			return cmds, err
		}
	}

	// prepare result
	for _, cmd := range commands {
		cmds.Command = append(cmds.Command, cmd.Command)
	}
	return cmds, nil
}

func getCustomCommandsByBuilder(customCommands Builders, builder string) ([]Commands, string, error) {
	for _, b := range customCommands.Builders {
		if builder == b.Name {
			return b.Type, b.Info, nil
		}
	}

	return nil, "", fmt.Errorf("the %s builder is not defined in the custom commands configuration", builder)

}

// CmdConverter - path and commands to execute
func CmdConverter(mPath string, cmdList []string) [][]string {
	var cmd [][]string
	for i := 0; i < len(cmdList); i++ {
		cmd = append(cmd, append([]string{mPath}, strings.Split(cmdList[i], " ")...))
	}
	return cmd
}

// GetModuleAndCommands - Get module from mta.yaml and
// commands (with resolved paths) configured for the module type
func GetModuleAndCommands(loc dir.IMtaParser, module string) (*mta.Module, []string, error) {
	mtaObj, err := loc.ParseFile()
	if err != nil {
		return nil, nil, err
	}
	// Get module respective command's to execute
	return moduleCmd(mtaObj, module)
}

// Get commands for specific module type
func moduleCmd(mta *mta.MTA, moduleName string) (*mta.Module, []string, error) {
	for _, m := range mta.Modules {
		if m.Name == moduleName {
			commandProvider, err := CommandProvider(*m)
			if err != nil {
				return nil, nil, err
			}
			return m, commandProvider.Command, nil
		}
	}
	return nil, nil, errors.Errorf("the %v module is not defined in the .mta file", moduleName)
}
