package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta/mta"
)

const (
	builderParam  = "builder"
	optionsSuffix = "-opts"
)

// CommandList - list of command to execute
type CommandList struct {
	Info    string
	Command []string
}

// GetBuilder - gets builder type of the module and indicator of custom builder
func GetBuilder(module *mta.Module) (string, bool, map[string]string) {
	// builder defined in build params is prioritised
	builderName := module.Type
	customBuilder := false
	if module.BuildParams != nil {
		if module.BuildParams[builderParam] != nil {
			builderName = module.BuildParams[builderParam].(string)
			customBuilder = true
		}
		options := getOpts(module)
		return builderName, customBuilder, options
	}
	// default builder is defined by type property of the module
	return module.Type, customBuilder, nil
}

// Get options for builder from mta.yaml
func getOpts(module *mta.Module) map[string]string {
	options := make(map[string]string)
	for paramName, paramValue := range module.BuildParams {
		if strings.HasSuffix(paramName, optionsSuffix) {
			optsName := strings.TrimSuffix(paramName, optionsSuffix)
			paramMapI, ok := paramValue.(map[interface{}]interface{})
			if ok {
				paramMapOpts := convertMapI(paramMapI, optsName)
				for paramMapKey, paramMapValue := range paramMapOpts {
					options[paramMapKey] = paramMapValue
				}
			} else {
				paramMap, ok := paramValue.(map[string]interface{})
				if ok {
					paramMapOpts := convertMap(paramMap, optsName)
					for paramMapKey, paramMapValue := range paramMapOpts {
						options[paramMapKey] = paramMapValue
					}
				} else {
					// options as sequence => jon them separated by space
					paramSlice, ok := paramValue.([]string)
					if ok {
						options[optsName] = strings.Join(paramSlice, " ")
					}
				}
			}
		}
	}

	return options
}

// Convert type map[interface{}]interface{} to map[string]string
func convertMapI(m map[interface{}]interface{}, optsName string) map[string]string {
	res := make(map[string]string)
	for key, value := range m {
		strKey, strValue := convertKeyValue(key.(string), value, optsName)
		res[strKey] = strValue
	}

	return res
}

// Convert type map[string]interface{} to map[string]string
func convertMap(m map[string]interface{}, optsName string) map[string]string {
	res := make(map[string]string)
	for key, value := range m {
		strKey, strValue := convertKeyValue(key, value, optsName)
		res[strKey] = strValue
	}
	return res
}

func convertKeyValue(key string, value interface{}, prefix string) (string, string) {
	strValue, ok := value.(string)
	if !ok {
		mapValueI := value.(map[interface{}]interface{})
		mapValue := convertMapI(mapValueI, "")
		strValue = ""
		for deepKey, deepValue := range mapValue {
			strValue = strValue + " --" + deepKey + " " + deepValue
		}
	}

	if prefix == "" {
		return key, strValue
	}
	return prefix + "." + key, strValue
}

// CommandProvider - Get build command's to execute
//noinspection GoExportedFuncWithUnexportedType
func CommandProvider(module mta.Module, commands []string) (CommandList, string, error) {
	// Get config from ./commands_cfg.yaml as generated artifacts from source
	moduleTypes, err := parseModuleTypes(ModuleTypeConfig)
	if err != nil {
		return CommandList{}, "", errors.Wrap(err, "failed to parse the module types configuration")
	}
	builderTypes, err := parseBuilders(BuilderTypeConfig)
	if err != nil {
		return CommandList{}, "", errors.Wrap(err, "failed to parse the builder types configuration")
	}
	return mesh(&module, &moduleTypes, builderTypes, commands)
}

// Match the object according to type and provide the respective command
func mesh(module *mta.Module, moduleTypes *ModuleTypes, builderTypes Builders, exec []string) (CommandList, string, error) {
	// The object support deep struct for future use, can be simplified to flat object
	var cmds CommandList
	var commands []Command
	var err error

	// get builder - module type name or custom builder if defined
	// and indicator if custom builder
	builder, custom, options := GetBuilder(module)

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

	if custom && builder != "_dummyBuilder_" {
		// custom builder used => get commands and info
		commands, cmds.Info, buildResults, err = getCustomCommandsByBuilder(builderTypes, builder)
		if err != nil {
			return cmds, "", err
		}
	}

	commands = addExecCommands(commands, exec)

	// prepare result
	cmds, buildResults = prepareMeshResult(cmds, buildResults, commands, options)
	return cmds, buildResults, nil
}

func addExecCommands(commands []Command, exec []string) []Command {
	result := commands
	if exec != nil {
		for _, cmd := range exec {
			result = append(result, Command{cmd})
		}
	}
	return result
}

// prepare commands list - mesh result
func prepareMeshResult(cmds CommandList, buildResults string, commands []Command, options map[string]string) (CommandList, string) {
	reg := regexp.MustCompile("{{\\w+([.]\\w+){0,1}}}")
	for _, cmd := range commands {
		if options != nil {
			cmd.Command = meshOpts(cmd.Command, options)
		}
		cmd.Command = reg.ReplaceAllString(cmd.Command, "")
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

func getCustomCommandsByBuilder(customCommands Builders, builder string) ([]Command, string, string, error) {
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
			commandProvider, buildResults, err := CommandProvider(*m, nil)
			if err != nil {
				return nil, nil, "", err
			}
			return m, commandProvider.Command, buildResults, nil
		}
	}
	return nil, nil, "", errors.Errorf(`the "%v" module is not defined in the MTA file`, moduleName)
}
