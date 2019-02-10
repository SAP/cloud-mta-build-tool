package buildops

import (
	"path/filepath"
	"reflect"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/fs"
	"github.com/SAP/cloud-mta/mta"
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

// GetBuilder - gets builder type of the module and indicator of custom builder
func GetBuilder(module *mta.Module) (string, bool) {
	// builder defined in build params is prioritised
	if module.BuildParams != nil && module.BuildParams[builderParam] != nil {
		return module.BuildParams[builderParam].(string), true
	}
	// default builder is defined by type property of the module
	return module.Type, false
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
			"processing requirements of the %v module based on the %v module failed when getting the %v module",
			moduleName, requires.Name, moduleName)
	}

	requiredModule, err := mta.GetModuleByName(requires.Name)
	if err != nil {
		return errors.Wrapf(err,
			"processing requirements of the %v module based on the %v module failed when getting the %v module",
			moduleName, requires.Name, requires.Name)
	}

	// Build paths for artifacts copying
	sourcePath := GetBuildResultsPath(ep, requiredModule)
	targetPath := getRequiredTargetPath(ep, module, requires)

	// execute copy of artifacts
	err = dir.CopyByPatterns(sourcePath, targetPath, requires.Artifacts)
	if err != nil {
		return errors.Wrapf(err,
			"processing requirements of the %v module based on the %v module failed when copying artifacts",
			moduleName, requiredModule.Name)
	}
	return nil
}

// GetBuildResultsPath - provides path of build results
func GetBuildResultsPath(ep dir.ISourceModule, module *mta.Module) string {
	path := ep.GetSourceModuleDir(module.Path)

	// if no sub-folder provided - build results will be saved in the module folder
	if module.BuildParams != nil && module.BuildParams[buildResultParam] != nil {
		// if sub-folder provided - build results are located in the subfolder of the module folder
		path = filepath.Join(path, module.BuildParams[buildResultParam].(string))
	}
	return path
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
			if p == platform {
				return true
			}
		}
		return false
	}
	sp := supportedPlatforms.([]interface{})
	for _, p := range sp {
		if p.(string) == platform {
			return true
		}
	}
	return false
}
