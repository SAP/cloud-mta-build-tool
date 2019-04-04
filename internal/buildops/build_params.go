package buildops

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta/mta"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
)

const (
	// SupportedPlatformsParam - name of build-params property for supported platforms
	SupportedPlatformsParam = "supported-platforms"
	builderParam            = "builder"
	requiresParam           = "requires"
	buildResultParam        = "build-result"
	nameParam               = "name"
	artifactsParam          = "artifacts"
	targetPathParam         = "target-path"
)

// BuildRequires - build requires section.
type BuildRequires struct {
	Name       string   `yaml:"name,omitempty"`
	Artifacts  []string `yaml:"artifacts,omitempty"`
	TargetPath string   `yaml:"target-path,omitempty"`
}

// getBuildRequires - gets Requires property of module's build-params property
// as generic property and converts it to slice of BuildRequires structures
func getBuildRequires(module *mta.Module) []BuildRequires {
	// check existence of module's build-params.require property
	if module.BuildParams != nil && module.BuildParams[requiresParam] != nil {
		requires := module.BuildParams[requiresParam].([]interface{})
		buildRequires := []BuildRequires{}
		// go through requirements
		for _, reqI := range requires {
			// cast requirement to generic map
			reqMap := reqI.(map[interface{}]interface{})
			// init resulting typed requirement
			reqStr := BuildRequires{
				Name:       getStrParam(reqMap, nameParam),
				Artifacts:  []string{},
				TargetPath: getStrParam(reqMap, targetPathParam),
			}
			// fill Artifacts field of resulting requirement
			if reqMap[artifactsParam] == nil {
				reqStr.Artifacts = nil
			} else {
				for _, artifact := range reqMap[artifactsParam].([]interface{}) {
					reqStr.Artifacts = append(reqStr.Artifacts, []string{artifact.(string)}...)
				}
			}
			// add typed requirement to result
			buildRequires = append(buildRequires, []BuildRequires{reqStr}...)

		}
		return buildRequires
	}
	return nil
}

// getStrParam - get string parameter from the generic map
func getStrParam(m map[interface{}]interface{}, param string) string {
	if m[param] == nil {
		return ""
	}
	return m[param].(string)
}

// Order of modules building is done according to the dependencies defined in build parameters.
// In case of problems in this definition build process should not start and corresponding error must be provided.
// Possible problems:
// 1.	Cyclic dependencies
// 2.	Dependency on not defined module

// ProcessRequirements - Processes build requirement of module (using moduleName).
func ProcessRequirements(ep dir.ISourceModule, mta *mta.MTA, requires *BuildRequires, moduleName string) error {

	// validate module names - both in process and required
	module, err := mta.GetModuleByName(moduleName)
	if err != nil {
		return errors.Wrapf(err,
			`processing requirements of the "%v" module based on the "%v" module failed when getting the "%v" module`,
			moduleName, requires.Name, moduleName)
	}

	requiredModule, err := mta.GetModuleByName(requires.Name)
	if err != nil {
		return errors.Wrapf(err,
			`processing requirements of the "%v" module based on the "%v" module failed when getting the "%v" module`,
			moduleName, requires.Name, requires.Name)
	}

	_, defaultBuildResult, err := commands.CommandProvider(*requiredModule)
	if err != nil {
		return errors.Wrapf(err,
			`processing requirements of the "%v" module based on the "%v" module failed when getting the "%v" module commands`,
			moduleName, requires.Name, requires.Name)
	}

	// Build paths for artifacts copying
	sourcePath, _, err := GetBuildResultsPath(ep, requiredModule, defaultBuildResult)
	if err != nil {
		return errors.Wrapf(err,
			`processing requirements of the "%v" module based on the "%v" module failed when getting the build results path`,
			moduleName, requires.Name)
	}
	targetPath := getRequiredTargetPath(ep, module, requires)

	// execute copy of artifacts
	err = dir.CopyByPatterns(sourcePath, targetPath, requires.Artifacts)
	if err != nil {
		return errors.Wrapf(err,
			`processing requirements of the "%v" module based on the "%v" module failed when copying artifacts`,
			moduleName, requiredModule.Name)
	}
	return nil
}

// GetBuildResultsPath - provides path of build results
func GetBuildResultsPath(ep dir.ISourceModule, module *mta.Module, defaultBuildResult string) (string, bool, error) {
	var path string
	if module.Path != "" {
		path = ep.GetSourceModuleDir(module.Path)
	} else {
		return "", false, nil
	}

	buildResultsDefined := false
	// if no sub-folder provided - build results will be saved in the module folder
	if module.BuildParams != nil && module.BuildParams[buildResultParam] != nil {
		// if sub-folder provided - build results are located in the subfolder of the module folder
		path = filepath.Join(path, module.BuildParams[buildResultParam].(string))
		buildResultsDefined = true
	} else if defaultBuildResult != "" {
		path = filepath.Join(path, defaultBuildResult)
		buildResultsDefined = true
	}

	if buildResultsDefined {
		sourceEntries, err := filepath.Glob(path)
		if err != nil {
			return "", false, err
		} else if len(sourceEntries) == 0 {
			return "", false, fmt.Errorf(`no entry found that matches the "%s" build results`, path)
		}
		return sourceEntries[0], true, nil
	}
	return path, false, nil
}

// getRequiredTargetPath - provides path of required artifacts
func getRequiredTargetPath(ep dir.ISourceModule, module *mta.Module, requires *BuildRequires) string {
	path := ep.GetSourceModuleDir(module.Path)
	if requires.TargetPath != "" {
		// if target folder provided - artifacts will be saved in the sub-folder of the module folder
		path = filepath.Join(path, requires.TargetPath)
	}
	return path
}

// PlatformDefined - if platform defined
// If platforms parameter not defined then no limitations on platform, method returns true
// Non empty list of platforms has to contain specific platform
func PlatformDefined(module *mta.Module, platform string) bool {
	if module.BuildParams == nil || module.BuildParams[SupportedPlatformsParam] == nil {
		return true
	}
	supportedPlatforms := module.BuildParams[SupportedPlatformsParam]
	if reflect.TypeOf(supportedPlatforms).Elem().Kind() == reflect.String {
		sp := supportedPlatforms.([]string)
		for _, p := range sp {
			if strings.ToLower(p) == platform {
				return true
			}
		}
		return false
	}
	sp := supportedPlatforms.([]interface{})
	for _, p := range sp {
		if strings.ToLower(p.(string)) == platform {
			return true
		}
	}
	return false
}
