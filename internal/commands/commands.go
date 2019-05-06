package commands

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta/mta"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

const (
	builderParam  = "builder"
	commandsParam = "commands"
	customBuilder = "custom"
	optionsSuffix = "-opts"
)

// CommandList - list of command to execute
type CommandList struct {
	Info    string
	Command []string
}

// GetBuilder - gets builder type of the module and indicator of custom builder
func GetBuilder(module *mta.Module) (string, bool, map[string]string, []string, error) {
	// builder defined in build params is prioritised
	if module.BuildParams != nil && module.BuildParams[builderParam] != nil {
		builderName := module.BuildParams[builderParam].(string)
		optsParamName := builderName + optionsSuffix
		// get options for builder from mta.yaml
		options := getOpts(module, optsParamName)
		var cmds []string
		if builderName == customBuilder {
			cmdsParam, ok := module.BuildParams[commandsParam]
			if !ok {
				logs.Logger.Warn(`the "commands" property is missing in the "custom" builder`)
				return builderName, true, options, []string{}, nil
			}
			cmds, ok = cmdsParam.([]string)
			if !ok {
				cmdsI, okI := cmdsParam.([]interface{})
				if okI {
					ok = true
					for _, cmdI := range cmdsI{
						cmd, okCmd := cmdI.(string)
						if !okCmd {
							ok=false
							break
						}
						cmds = append(cmds, cmd)
					}
				}
			}
			if !ok {
				return builderName, true, options, cmds, fmt.Errorf(`the "commands" property is defined incorrectly; the property must contain a sequence of strings`)
			}
		}

		return builderName, true, options, cmds, nil
	}
	// default builder is defined by type property of the module
	return module.Type, false, nil, nil, nil
}

// Get options for builder from mta.yaml
func getOpts(module *mta.Module, optsParamName string) map[string]string {
	options := module.BuildParams[optsParamName]
	optionsMap := make(map[string]string)
	if options != nil {
		optionsMap = convert(options.(map[interface{}]interface{}))
	}

	return optionsMap
}

// Convert type map[interface{}]interface{} to map[string]string
func convert(m map[interface{}]interface{}) map[string]string {
	res := make(map[string]string)
	for key, value := range m {
		strKey := key.(string)
		strValue := value.(string)

		res[strKey] = strValue
	}

	return res
}

// CommandProvider - Get build command's to execute
//noinspection GoExportedFuncWithUnexportedType
func CommandProvider(module mta.Module) (CommandList, string, error) {
	// Get config from ./commands_cfg.yaml as generated artifacts from source
	moduleTypes, err := parseModuleTypes(ModuleTypeConfig)
	if err != nil {
		return CommandList{}, "", errors.Wrap(err, "failed to parse the module types configuration")
	}
	builderTypes, err := parseBuilders(BuilderTypeConfig)
	if err != nil {
		return CommandList{}, "", errors.Wrap(err, "failed to parse the builder types configuration")
	}
	return mesh(&module, &moduleTypes, &builderTypes)
}

// Match the object according to type and provide the respective command
func mesh(module *mta.Module, moduleTypes *ModuleTypes, builderTypes *Builders) (CommandList, string, error) {
	// The object support deep struct for future use, can be simplified to flat object
	var cmds CommandList
	var cmdList []string
	var commands []Command
	var err error

	// get builder - module type name or custom builder if defined
	// and indicator if custom builder
	builder, custom, options, cmdList, err := GetBuilder(module)
	if err != nil {
		return CommandList{Command: []string{}}, "", err
	}

	// if module type used - get from module types configuration corresponding commands or custom builder if defined
	if !custom {
		for _, m := range moduleTypes.ModuleTypes {
			if m.Name == builder {
				if m.Builder != "" {
					// custom builder defined
					// check that no commands defined for module type
					if m.Commands != nil && len(m.Commands) > 0 {
						return cmds, "", fmt.Errorf(
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

	buildResults := ""

	if custom {
		// custom builder used => get commands and info
		commands, cmds.Info, buildResults, err = getCustomCommandsByBuilder(builderTypes, builder, cmdList)
		if err != nil {
			return cmds, "", err
		}
	}

	// prepare result
	cmds, buildResults = prepareMeshResult(cmds, buildResults, commands, options)
	return cmds, buildResults, nil
}

// prepare commands list - mesh result
func prepareMeshResult(cmds CommandList, buildResults string, commands []Command, options map[string]string) (CommandList, string) {
	for _, cmd := range commands {
		if options != nil {
			cmd.Command = meshOpts(cmd.Command, options)
		}
		cmds.Command = append(cmds.Command, cmd.Command)
	}
	return cmds, buildResults
}

// Update command according to options arguments
func meshOpts(cmd string, options map[string]string) string {
	c := cmd
	for key, value := range options {
		c = strings.Replace(c, "{{"+key+"}}", value, -1)
	}
	return c
}

func getCustomCommandsByBuilder(customCommands *Builders, builder string, cmds []string) ([]Command, string, string, error) {
	if builder == customBuilder {
		var res []Command
		for _, cmd := range cmds {
			res = append(res, Command{cmd})
		}
		return res, "", "", nil
	}

	for _, b := range customCommands.Builders {
		if builder == b.Name {
			return b.Commands, b.Info, b.BuildResult, nil
		}
	}

	return nil, "", "", fmt.Errorf(`the "%s" builder is not defined in the custom commands configuration`, builder)

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
func GetModuleAndCommands(loc dir.IMtaParser, module string) (*mta.Module, []string, string, error) {
	mtaObj, err := loc.ParseFile()
	if err != nil {
		return nil, nil, "", err
	}
	// Get module respective command's to execute
	return moduleCmd(mtaObj, module)
}

// Get commands for specific module type
func moduleCmd(mta *mta.MTA, moduleName string) (*mta.Module, []string, string, error) {
	for _, m := range mta.Modules {
		if m.Name == moduleName {
			commandProvider, buildResults, err := CommandProvider(*m)
			if err != nil {
				return nil, nil, "", err
			}
			return m, commandProvider.Command, buildResults, nil
		}
	}
	return nil, nil, "", errors.Errorf(`the "%v" module is not defined in the MTA file`, moduleName)
}
