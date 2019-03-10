package commands

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta/mta"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/buildops"
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
	moduleTypes, err := parseModuleTypes(ModuleTypeConfig)
	if err != nil {
		return CommandList{}, errors.Wrap(err, "failed to parse the module types configuration")
	}
	builderTypes, err := parseBuilders(BuilderTypeConfig)
	if err != nil {
		return CommandList{}, errors.Wrap(err, "failed to parse the builder types configuration")
	}
	return mesh(&modules, &moduleTypes, builderTypes)
}

// Match the object according to type and provide the respective command
func mesh(module *mta.Module, moduleTypes *ModuleTypes, builderTypes Builders) (CommandList, error) {
	// The object support deep struct for future use, can be simplified to flat object
	var cmds CommandList
	var commands []Command
	var err error

	// get builder - module type name or custom builder if defined
	// indicator if custom builder
	// options of builder if defined
	builder, custom, options := buildops.GetBuilder(module)

	// if module type used - get from module types configuration corresponding commands or custom builder if defined
	if !custom {
		for _, m := range moduleTypes.ModuleTypes {
			if m.Name == builder {
				if m.Builder != "" {
					// custom builder defined
					// check that no commands defined for module type
					if m.Commands != nil && len(m.Commands) > 0 {
						return cmds, fmt.Errorf(
							"the module type definition can include either the builder or the commands; the %s module type includes both",
							m.Name)
					}
					// continue with custom builders search
					builder = m.Builder
					custom = true
				} else {
					// get related information
					cmds.Info = m.Info
					commands = m.Commands
				}
			}
		}
	}

	if custom {
		// custom builder used => get commands and info
		commands, cmds.Info, err = getCustomCommandsByBuilder(builderTypes, builder, options)
		if err != nil {
			return cmds, err
		}
	}

	// prepare result
	for _, cmd := range commands {
		if options != nil {
			cmd.Command = meshOpts(cmd.Command, options)
		}
		cmds.Command = append(cmds.Command, cmd.Command) //change cmd.Command to c
	}
	return cmds, nil
}

func meshOpts(cmd string, options map[string]string) string {
	c := cmd
	for key, value := range options {
		r := strings.NewReplacer("{{"+key+"}}", value)
		c = r.Replace(c)
	}
	return c
}

func getCustomCommandsByBuilder(customCommands Builders, builder string, options map[string]string) ([]Command, string, error) {
	for _, b := range customCommands.Builders {
		if builder == b.Name {
			if b.BuilderTypes != nil {
				return getCustomCommandsByBuilderType(b, options)
			}

			return b.Commands, b.Info, nil
		}
	}

	return nil, "", fmt.Errorf(`the "%s" builder is not defined in the custom commands configuration`, builder)
}

func getCustomCommandsByBuilderType(customCommands builder, options map[string]string) ([]Command, string, error) {
	for _, b := range customCommands.BuilderTypes {
		if options["repo-type"] == b.Name {
			return b.Commands, b.Info, nil
		}
	}

	return nil, "", fmt.Errorf(`the "%s" builder type is not defined in the custom commands configuration`, options["repo-type"])
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
	return nil, nil, errors.Errorf(`the "%v" module is not defined in the MTA file`, moduleName)
}
