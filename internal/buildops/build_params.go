package buildops

import (
	"path/filepath"

	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/mta"
)

const (
	builderParam            = "builder"
	SupportedPlatformsParam = "supported-platforms"
	requiresParam           = "requires"
	buildResultsParam       = "build-results"
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

func GetBuilder(module *mta.Module) string {
	if module.BuildParams != nil && module.BuildParams[builderParam] != nil {
		return module.BuildParams[builderParam].(string)
	}
	return module.Type
}

func getRequires(module *mta.Module) []BuildRequires {
	if module.BuildParams != nil && module.BuildParams[requiresParam] != nil {
		requires := module.BuildParams[requiresParam].([]interface{})
		buildRequires := []BuildRequires{}
		for _, reqI := range requires {
			reqMap := reqI.(map[interface{}]interface{})
			reqStr := BuildRequires{
				Name:       getStrParam(reqMap, nameParam),
				Artifacts:  []string{},
				TargetPath: getStrParam(reqMap, targetPathParam),
			}
			if reqMap[artifactsParam] == nil {
				reqStr.Artifacts = nil
			} else {
				for _, artifact := range reqMap[artifactsParam].([]interface{}) {
					reqStr.Artifacts = append(reqStr.Artifacts, []string{artifact.(string)}...)
				}
			}
			buildRequires = append(buildRequires, []BuildRequires{reqStr}...)

		}
		return buildRequires
	}
	return nil
}

func getStrParam(m map[interface{}]interface{}, param string) string {
	if m[param] == nil {
		return ""
	} else {
		return m[param].(string)
	}
}

// Order of modules building is done according to the dependencies defined in build parameters.
// In case of problems in this definition build process should not start and corresponding error must be provided.
// Possible problems:
// 1.	Cyclic dependencies
// 2.	Dependency on not defined module

// ProcessRequirements - Processes build requirement of module (using moduleName).
func ProcessRequirements(ep *dir.Loc, mta *mta.MTA, requires *BuildRequires, moduleName string) error {

	// validate module names - both in process and required
	module, err := mta.GetModuleByName(moduleName)
	if err != nil {
		return errors.Wrapf(err, "Processing requirements of module %v based on module %v failed on getting module", moduleName, requires.Name)
	}
	requiredModule, err := mta.GetModuleByName(requires.Name)
	if err != nil {
		return errors.Wrapf(err, "Processing requirements of module %v based on module %v failed on getting required module", moduleName, requires.Name)
	}
	// Get slice of artifacts
	artifacts := requires.Artifacts

	// Build paths for artifacts copying
	sourcePath, err := getBuildResultsPath(ep, requiredModule)
	if err != nil {
		return errors.Wrapf(err, "Processing requirements of module %v based on module %v failed on getting Results Path", moduleName, requiredModule.Name)
	}
	targetPath, err := getRequiredTargetPath(ep, module, requires)
	if err != nil {
		return errors.Wrapf(err, "Processing requirements of module %v based on module %v failed on getting Required Target Path", moduleName, requiredModule.Name)
	}
	// execute copy of artifacts
	err = dir.CopyByPatterns(sourcePath, targetPath, artifacts)

	if err != nil {
		return errors.Wrapf(err, "Processing requirements of module %v based on module %v failed on artifacts copying", moduleName, requiredModule.Name)
	}
	return nil
}

// getBuildResultsPath - provides path of build results
func getBuildResultsPath(ep *dir.Loc, module *mta.Module) (string, error) {
	path, err := ep.GetSourceModuleDir(module.Path)
	if err != nil {
		return "", errors.Wrapf(err, "getBuildResultsPath failed getting directory of module %v", module.Path)
	}
	// if no sub-folder provided - build results will be saved in the module folder
	if module.BuildParams != nil && module.BuildParams[buildResultsParam] != nil {
		// if sub-folder provided - build results are located in the subfolder of the module folder
		path = filepath.Join(path, module.BuildParams[buildResultsParam].(string))
	}
	return path, nil
}

// getRequiredTargetPath - provides path of required artifacts
func getRequiredTargetPath(ep *dir.Loc, module *mta.Module, requires *BuildRequires) (string, error) {
	path, err := ep.GetSourceModuleDir(module.Path)
	if err != nil {
		return "", errors.Wrapf(err, "getRequiredTargetPath failed getting directory of module %v", module.Name)
	}
	if requires.TargetPath != "" {
		// if target folder provided - artifacts will be saved in the sub-folder of the module folder
		path = filepath.Join(path, requires.TargetPath)
	}
	return path, nil
}

// PlatformsDefined - if platforms defined
// Only empty list of platforms indicates no platforms defined
func PlatformsDefined(module *mta.Module) bool {
	if module.BuildParams == nil || module.BuildParams[SupportedPlatformsParam] == nil {
		return true
	}
	supportedPlatforms := module.BuildParams[SupportedPlatformsParam].([]string)
	return len(supportedPlatforms) > 0
}
