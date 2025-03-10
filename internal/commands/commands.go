package commands

import (
	"fmt"
	"strings"

	"github.com/kballard/go-shellquote"
	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta/mta"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

const (
	builderParam                 = "builder"
	commandsParam                = "commands"
	customBuilder                = "custom"
	golangBuilder                = "golang"
	optionsSuffix                = "-opts"
	goModuleType                 = "go"
	cyclonedx_npm                = "@cyclonedx/cyclonedx-npm"
	cyclonedx_npm_version        = "1.19.3"
	cyclonedx_npm_schema_version = "1.4"
)

// CommandList - list of command to execute
type CommandList struct {
	Info    string
	Command []string
}

// GetBuilder - gets builder type of the module and indicator of custom builder
// if build-parameter == null or build-parameter.builder == null, return builder=module.type and custom=false
// else if build-paramete.builder != custom, return builder=build-paramete.builder and custom=true
// else if build-paramete.builder == custom, return builder=custom and custom=true
func GetBuilder(module *mta.Module) (string, bool, map[string]string, []string, error) {
	// builder defined in build params is prioritised
	if module.BuildParams != nil && module.BuildParams[builderParam] != nil {
		builderName := module.BuildParams[builderParam].(string)
		checkDeprecatedBuilder(builderName)
		optsParamName := builderName + optionsSuffix
		// get options for builder from mta.yaml
		options := getOpts(module, optsParamName)
		var cmds []string
		if builderName == customBuilder {
			cmdsParam, ok := module.BuildParams[commandsParam]
			if !ok {
				logs.Logger.Warn(missingPropMsg)
				return builderName, true, options, []string{}, nil
			}
			cmds, ok = cmdsParam.([]string)
			if !ok {
				cmdsI, okI := cmdsParam.([]interface{})
				if okI {
					ok = true
					for _, cmdI := range cmdsI {
						cmd, okCmd := cmdI.(string)
						if !okCmd {
							ok = false
							break
						}
						cmds = append(cmds, cmd)
					}
				}
			}
			if !ok {
				return builderName, true, options, cmds, fmt.Errorf(wrongPropMsg)
			}
		}

		return builderName, true, options, cmds, nil
	}
	// default builder is defined by type property of the module
	return module.Type, false, nil, nil, nil
}

func isNativeBuilderType(builderName string) (bool, error) {
	builderTypes, err := parseBuilders(BuilderTypeConfig)
	if err != nil {
		return false, errors.Wrap(err, parseBuilderCfgFailedMsg)
	}

	for _, b := range builderTypes.Builders {
		if builderName == b.Name {
			return true, err
		}
	}
	return false, err
}

func getSBomBuilderByModuleType(typeName string) (bool, string, error) {
	moduleTypes, err := parseModuleTypes(ModuleTypeConfig)
	if err != nil {
		return false, "", errors.Wrap(err, parseModuleCfgFailedMsg)
	}

	for _, t := range moduleTypes.ModuleTypes {
		if typeName == t.Name {
			return true, t.Builder, nil
		}
	}

	// if module.type == go, return the golang builder;
	// Notice, there is no go type in ModuleTypeConfig(module_type_cfg.yaml)
	if typeName == goModuleType {
		return true, golangBuilder, nil
	}

	return false, "", errors.Wrapf(err, notNativeModuleTypeMsg, typeName)
}

func getModuleSBomBuilder(module *mta.Module) (string, error) {
	var builderName string
	var err error

	// get builder by build-parameter.builder
	if module.BuildParams != nil && module.BuildParams[builderParam] != nil {
		builderName = module.BuildParams[builderParam].(string)
		checkDeprecatedBuilder(builderName)
		if builderName == customBuilder {
			_, ok := module.BuildParams[commandsParam]
			if !ok {
				return builderName, errors.Wrap(err, missingPropMsg)
			}
			return builderName, nil
		}

		// check if builder is native builder (builder_type_cfg.yaml)
		isnativebuilder, err := isNativeBuilderType(builderName)
		if !isnativebuilder {
			return builderName, errors.Wrapf(err, notNativeBuilderMsg, builderName)
		}

		return builderName, nil
	}

	// get builder by module type
	isfind, builderName, err := getSBomBuilderByModuleType(module.Type)
	if !isfind {
		return "", err
	}
	return builderName, nil
}

// Get options for builder from mta.yaml
func getOpts(module *mta.Module, optsParamName string) map[string]string {
	options := module.BuildParams[optsParamName]
	optionsMap := make(map[string]string)
	if options != nil {
		optionsMapS, ok := options.(map[string]interface{})
		if !ok {
			optionsMapS = ConvertMap(options.(map[interface{}]interface{}))
		}
		optionsMap = convert(optionsMapS)
	}

	return optionsMap
}

// Convert type map[string]interface{} to map[string]string
func convert(m map[string]interface{}) map[string]string {
	res := make(map[string]string)
	for strKey, value := range m {
		strValue := ""
		if value != nil {
			strValue = value.(string)
		}
		res[strKey] = strValue
	}

	return res
}

// ConvertMap converts type map[interface{}]interface{} to map[string]interface{}
func ConvertMap(m map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for key, value := range m {
		strKey := key.(string)
		res[strKey] = value
	}

	return res
}

// CommandProvider - Get build command's to execute
// noinspection GoExportedFuncWithUnexportedType
func CommandProvider(module mta.Module) (CommandList, string, error) {
	// Get config from ./commands_cfg.yaml as generated artifacts from source
	moduleTypes, err := parseModuleTypes(ModuleTypeConfig)
	if err != nil {
		return CommandList{}, "", errors.Wrap(err, parseModuleCfgFailedMsg)
	}
	builderTypes, err := parseBuilders(BuilderTypeConfig)
	if err != nil {
		return CommandList{}, "", errors.Wrap(err, parseBuilderCfgFailedMsg)
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
						return cmds, "", fmt.Errorf(wrongModuleTypeDefMsg, m.Name)
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

	return nil, "", "", fmt.Errorf(undefinedBuilderMsg, builder)
}

// CmdConverter - path and commands to execute
func CmdConverter(mPath string, cmdList []string) ([][]string, error) {
	var cmd [][]string
	for i := 0; i < len(cmdList); i++ {
		split, err := shellquote.Split(cmdList[i])
		if err != nil {
			return nil, errors.Wrapf(err, BadCommandMsg, cmdList[i])
		}
		cmd = append(cmd, append([]string{mPath}, split...))
	}
	return cmd, nil
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
	return nil, nil, "", errors.Errorf(undefinedModuleMsg, moduleName)
}

// GetModuleSBomGenCommands - get sbom generate command for module
// if unknow sbom gen builder or custom builder, empty [][]string and nil error will be return
func GetModuleSBomGenCommands(loc *dir.Loc, module *mta.Module,
	sbomFileName string, sbomFileType string, sbomFileSuffix string) ([][]string, error) {
	var cmd string
	var cmds []string
	var commandList [][]string

	builder, err := getModuleSBomBuilder(module)
	if err != nil {
		return [][]string{}, err
	}

	switch builder {
	case "npm", "npm-ci", "grunt", "evo":
		cmd = "npm install"
		cmds = append(cmds, cmd)
		// cmd = "npm install " + cyclonedx_npm + "@" + cyclonedx_npm_version + " --no-save"
		// cmds = append(cmds, cmd)
		// cmd = "npx cyclonedx-npm --output-format " + strings.ToUpper(sbomFileType) + " --spec-version " + cyclonedx_npm_schema_version + " --output-file " + sbomFileName + sbomFileSuffix
		cmd = "npx " + cyclonedx_npm + "@" + cyclonedx_npm_version + " --output-format " + strings.ToUpper(sbomFileType) + " --spec-version " + cyclonedx_npm_schema_version + " --output-file " + sbomFileName + sbomFileSuffix
		cmds = append(cmds, cmd)
	case "golang":
		cmd = "cyclonedx-gomod mod -output-version 1.4 -licenses -output " + sbomFileName + sbomFileSuffix
		cmds = append(cmds, cmd)
	case "maven", "fetcher", "maven_deprecated":
		cmd = "mvn org.cyclonedx:cyclonedx-maven-plugin:2.9.0:makeAggregateBom " +
			"-DschemaVersion=1.4 -DincludeBomSerialNumber=true -DincludeCompileScope=true " +
			"-DincludeRuntimeScope=true -DincludeSystemScope=true -DincludeTestScope=false -DincludeLicenseText=false " +
			"-DoutputFormat=" + sbomFileType + " -DoutputName=" + sbomFileName + ".bom"
		cmds = append(cmds, cmd)
	case "custom":
		// first check if custom SBOM creation commands are provided
		customSbomGenCmds, ok := module.BuildParams["sbom-create-commands"].([]string)
		// in case no custom commands are provided use standard way of creating SBOM
		if !ok || (ok && len(customSbomGenCmds) == 0) {
			switch module.Type {
			case "nodejs":
				cmd = "npm install"
				cmds = append(cmds, cmd)
				cmd = "npx " + cyclonedx_npm + "@" + cyclonedx_npm_version + " --output-format " + strings.ToUpper(sbomFileType) + " --spec-version " + cyclonedx_npm_schema_version + " --output-file " + sbomFileName + sbomFileSuffix
				cmds = append(cmds, cmd)
			case "java":
				cmd = "mvn org.cyclonedx:cyclonedx-maven-plugin:2.7.5:makeAggregateBom " +
					"-DschemaVersion=1.4 -DincludeBomSerialNumber=true -DincludeCompileScope=true " +
					"-DincludeRuntimeScope=true -DincludeSystemScope=true -DincludeTestScope=false -DincludeLicenseText=false " +
					"-DoutputFormat=" + sbomFileType + " -DoutputName=" + sbomFileName + ".bom"
				cmds = append(cmds, cmd)
			}
			// in case custom SBOM creation commands are provided use them
		} else {
			// replace fileName placeholder ${sbom-file-name} which is to be provided in the custom SBOM creation commands
			for i := range customSbomGenCmds {
				customSbomGenCmds[i] = strings.ReplaceAll(customSbomGenCmds[i], "${sbom-file-name}", sbomFileName+sbomFileSuffix)
			}
			cmds = append(cmds, customSbomGenCmds...)
		}
	default:
	}

	modulePath := loc.GetSourceModuleDir(module.Path)
	commandList, err = CmdConverter(modulePath, cmds)
	if err != nil {
		return [][]string{}, err
	}
	return commandList, err
}

// GetSBomsMergeCommand - generate merge sbom file command under sbom tmp dir
// if empty sbomFileNames, return empty commandList, nil error
func GetSBomsMergeCommand(loc *dir.Loc, cyclonedx_cli string, mtaObj *mta.MTA, sbomTmpDir string, sbomFileNames []string,
	sbomName, sbomType, sbomSuffix string) ([][]string, error) {
	var cmd string
	var cmds []string
	var commandList [][]string

	// len(sbomFileName) should not be 0, if 0 then raise an error
	if len(sbomFileNames) == 0 {
		return commandList, errors.New(emptySBomFileInputMsg)
	}

	var inputFiles string
	for _, fileName := range sbomFileNames {
		inputFiles = inputFiles + " " + fileName + " "
	}

	// ./cyclonedx merge --input-files test_1.bom.xml test_2.bom.xml test_3.bom.xml --output-file merged.bom.xml
	cmd = cyclonedx_cli + " merge --input-files " + inputFiles + " --output-file " + sbomName +
		" --input-format " + sbomType + " --output-format " + sbomType + " --hierarchical" + " --name " + mtaObj.ID + " --version " + mtaObj.Version
	cmds = append(cmds, cmd)
	commandList, err := CmdConverter(sbomTmpDir, cmds)

	return commandList, err
}
