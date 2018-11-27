package builders

import (
	"cloud-mta-build-tool/mta"
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
	commands, err := parse(CommandsConfig)
	if err != nil {
		return CommandList{}, err
	}
	return mesh(modules, commands), nil
}

// Match the object according to type and provide the respective command
func mesh(modules mta.Module, commands Builders) CommandList {
	// The object support deep struct for future use, can be simplified to flat object
	var cmds CommandList
	for _, b := range commands.Builders {
		moduleType := modules.Type
		if modules.BuildParams.Builder != "" {
			moduleType = modules.BuildParams.Builder
		}
		// Return only matching types
		if moduleType == b.Name {
			cmds.Info = b.Info
			for _, cmd := range b.Type {
				cmds.Command = append(cmds.Command, cmd.Command)
			}
			break
		}
	}
	return cmds
}
