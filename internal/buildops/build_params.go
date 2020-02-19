package buildops

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta/mta"
)

const (
	// SupportedPlatformsParam - name of build-params property for supported platforms
	SupportedPlatformsParam = "supported-platforms"

	// ModuleArtifactDefaultName - the default name of the build artifact.
	// It can be changed using properties like build-result or build-artifact-name in the build parameters.
	ModuleArtifactDefaultName = "data.zip"
	builderParam              = "builder"
	requiresParam             = "requires"
	buildResultParam          = "build-result"
	nameParam                 = "name"
	artifactsParam            = "artifacts"
	buildArtifactNameParam    = "build-artifact-name"
	targetPathParam           = "target-path"
	noSourceParam             = "no-source"
)

// BuildRequires - build requires section.
type BuildRequires struct {
	Name       string   `yaml:"name,omitempty"`
	Artifacts  []string `yaml:"artifacts,omitempty"`
	TargetPath string   `yaml:"target-path,omitempty"`
}

// GetBuildRequires - gets Requires property of module's build-params property
// as generic property and converts it to slice of BuildRequires structures
func GetBuildRequires(module *mta.Module) []BuildRequires {
	// check existence of module's build-params.require property
	if module.BuildParams != nil && module.BuildParams[requiresParam] != nil {
		requires := module.BuildParams[requiresParam].([]interface{})
		buildRequires := []BuildRequires{}
		// go through requirements
		for _, reqI := range requires {
			// cast requirement to generic map
			reqMap, ok := reqI.(map[string]interface{})
			if !ok {
				reqMap = commands.ConvertMap(reqI.(map[interface{}]interface{}))
			}
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

// getStrParam - get string parameter from the map
func getStrParam(m map[string]interface{}, param string) string {
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

// GetRequiresArtifacts returns the source path, target path and patterns of files and folders to copy from a module's requires section
func GetRequiresArtifacts(ep dir.ISourceModule, mta *mta.MTA, requires *BuildRequires, moduleName string, resolveBuildResult bool) (source string, target string, patterns []string, err error) {
	// validate module names - both in process and required
	module, err := mta.GetModuleByName(moduleName)
	if err != nil {
		return "", "", nil, errors.Wrapf(err, reqFailedOnModuleGetMsg, moduleName, requires.Name, moduleName)
	}

	requiredModule, err := mta.GetModuleByName(requires.Name)
	if err != nil {
		return "", "", nil, errors.Wrapf(err, reqFailedOnModuleGetMsg, moduleName, requires.Name, requires.Name)
	}

	_, defaultBuildResult, err := commands.CommandProvider(*requiredModule)
	if err != nil {
		return "", "", nil, errors.Wrapf(err, reqFailedOnCommandsGetMsg, moduleName, requires.Name, requires.Name)
	}

	// Build paths for artifacts copying
	sourcePath, err := GetModuleSourceArtifactPath(ep, false, requiredModule, defaultBuildResult, resolveBuildResult)
	if err != nil {
		return "", "", nil, errors.Wrapf(err, reqFailedOnBuildResultMsg, moduleName, requires.Name)
	}
	targetPath := getRequiredTargetPath(ep, module, requires)
	return sourcePath, targetPath, requires.Artifacts, nil
}

// ProcessRequirements - Processes build requirement of module (using moduleName).
func ProcessRequirements(ep dir.ISourceModule, mta *mta.MTA, requires *BuildRequires, moduleName string) error {
	sourcePath, targetPath, artifacts, err := GetRequiresArtifacts(ep, mta, requires, moduleName, true)
	if err != nil {
		return err
	}
	// execute copy of artifacts
	err = dir.CopyByPatterns(sourcePath, targetPath, artifacts)
	if err != nil {
		return errors.Wrapf(err, reqFailedOnCopyMsg, moduleName, requires.Name)
	}
	return nil
}

// GetModuleSourceArtifactPath - get the module's artifact that has to be archived in the mtar, from the project sources
func GetModuleSourceArtifactPath(loc dir.ISourceModule, depDesc bool, module *mta.Module, defaultBuildResult string, resolveBuildResult bool) (path string, e error) {
	if module.Path == "" {
		return "", nil
	}
	path = loc.GetSourceModuleDir(module.Path)
	if !depDesc {
		buildResult := defaultBuildResult
		var ok bool
		if module.BuildParams != nil && module.BuildParams[buildResultParam] != nil {
			buildResult, ok = module.BuildParams[buildResultParam].(string)
			if !ok {
				return "", errors.Errorf(WrongBuildResultMsg, module.BuildParams[buildResultParam], module.Name)
			}
		}
		if buildResult != "" {
			path = filepath.Join(path, buildResult)
			if resolveBuildResult {
				path, e = dir.FindPath(path)
				if e != nil {
					return "", e
				}
			}
		}
	}
	return path, nil
}

// IsArchive - check if file is a folder or an archive
func IsArchive(path string) (isArchive bool, e error) {
	info, err := os.Stat(path)

	if err != nil {
		return false, err
	}
	isFolder := info.IsDir()
	isArchive = false
	if !isFolder {
		ext := filepath.Ext(path)
		isArchive = ext == ".zip" || ext == ".jar" || ext == ".war"
	}
	return isArchive, nil
}

// GetModuleTargetArtifactPath - get the path to where the module's artifact should be created in the temp folder, from which it's archived in the mtar
func GetModuleTargetArtifactPath(moduleLoc dir.IModule, depDesc bool, module *mta.Module, defaultBuildResult string) (path string, toArchive bool, e error) {
	if module.Path == "" {
		return "", false, nil
	}
	if depDesc {
		path = filepath.Join(moduleLoc.GetTargetModuleDir(module.Path))
	} else {
		moduleSourceArtifactPath, err := GetModuleSourceArtifactPath(moduleLoc, depDesc, module, defaultBuildResult, true)
		if err != nil {
			return "", false, err
		}
		isArchive, err := IsArchive(moduleSourceArtifactPath)
		if err != nil {
			return "", false, errors.Wrapf(err, wrongPathMsg, moduleSourceArtifactPath)
		}
		artifactName, artifactExt, err := getArtifactInfo(isArchive, module, moduleSourceArtifactPath)
		if err != nil {
			return "", false, err
		}
		toArchive = !isArchive

		artifactRelPath, err := moduleLoc.GetSourceModuleArtifactRelPath(module.Path, moduleSourceArtifactPath)
		if err != nil {
			return "", false, err
		}
		path = filepath.Join(moduleLoc.GetTargetModuleDir(module.Name), artifactRelPath, artifactName+artifactExt)
	}
	return path, toArchive, nil
}

func getArtifactInfo(isArchive bool, module *mta.Module, moduleSourceArtifactPath string) (artifactName, artifactExt string, err error) {
	var ok bool
	var artifactFullName string
	if isArchive {
		artifactFullName = filepath.Base(moduleSourceArtifactPath)
	} else {
		artifactFullName = ModuleArtifactDefaultName
	}
	artifactExt = filepath.Ext(artifactFullName)
	artifactName = artifactFullName[0 : len(artifactFullName)-len(artifactExt)]
	if module.BuildParams != nil && module.BuildParams[buildArtifactNameParam] != nil {
		artifactName, ok = module.BuildParams[buildArtifactNameParam].(string)
		if !ok {
			return "", "", errors.Errorf(WrongBuildArtifactNameMsg, module.BuildParams[buildArtifactNameParam], module.Name)
		}
	}
	return
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

// IfNoSource - checks if "no-source" build parameter defined and set to "true"
func IfNoSource(module *mta.Module) bool {
	if module.BuildParams != nil && module.BuildParams[noSourceParam] != nil {
		noSource, ok := module.BuildParams[noSourceParam].(bool)
		return ok && noSource
	}
	return false
}
