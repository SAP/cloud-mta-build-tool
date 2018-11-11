package builders

import (
	"cloud-mta-build-tool/mta"
)

// commandList - list of command to execute
type commandList struct {
	Info    string
	Command []string
}

// CommandProvider - Get build command's to execute
//noinspection GoExportedFuncWithUnexportedType
func CommandProvider(modules mta.Modules) (commandList, error) {
	// Get config from ./commands_cfg.yaml as generated artifacts from source
	commands, err := parse(CommandsConfig)
	if err != nil {
		return commandList{}, err
	}
	return mesh(modules, commands), nil
}

// Match the object according to type and provide the respective command
func mesh(modules mta.Modules, commands Builders) commandList {
	// The object support deep struct for future use, can be simplified to flat object
	var cmds commandList
	for _, b := range commands.Builders {
		// Return only matching types
		if modules.Type == b.Name {
			cmds.Info = b.Info
			for _, cmd := range b.Type {
				cmds.Command = append(cmds.Command, cmd.Command)
			}
			break
		}
	}
	return cmds
}
