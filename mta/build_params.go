package mta

import (
	"path/filepath"
	"strings"

	"cloud-mta-build-tool/internal/fsys"
	"github.com/pkg/errors"
)

// Order of modules building is done according to the dependencies defined in build parameters.
// In case of problems in this definition build process should not start and corresponding error must be provided.
// Possible problems:
// 1.	Cyclic dependencies
// 2.	Dependency on not defined module
// 3.	Wrong definition of artifacts

// ProcessRequirements - Processes build requirement of module (moduleName)
func (requires *BuildRequires) ProcessRequirements(ep *dir.MtaLocationParameters, mta *MTA, moduleName string) error {
	// validate module names - both in process and required
	module, err := mta.GetModuleByName(moduleName)
	if err != nil {
		return errors.Wrapf(err, "Processed module %v not defined in MTA", moduleName)
	}
	requiredModule, err := mta.GetModuleByName(requires.Name)
	if err != nil {
		return errors.Wrapf(err, "Required module %v not defined in MTA", requires.Name)
	}
	// Get slice of artifacts
	artifactsStr := strings.Replace(requires.Artifacts, "[", "", 1)
	artifactsStr = strings.Replace(artifactsStr, "]", "", 1)
	artifacts := strings.Split(artifactsStr, ",")

	// Validate artifacts
	err = validateArtifacts(ep, requiredModule, artifacts)
	if err != nil {
		return errors.Wrapf(err, "Error while processing requirements of module %v based on module %v", moduleName, requiredModule.Name)
	}
	// Build paths for artifacts copying
	sourcePath := requiredModule.getBuildResultsPath(ep)
	targetPath := requires.getRequiredTargetPath(ep, module)

	// execute copy of artifacts
	return CopyRequiredArtifacts(sourcePath, targetPath, artifacts)
}

// CopyRequiredArtifacts - copies artifacts of predecessor (source module) to dependent (target module)
func CopyRequiredArtifacts(sourcePath, targetPath string, artifacts []string) error {
	if len(artifacts) == 1 {
		if artifacts[0] == "*" {
			// copies all source module folder's entries
			dir.CopyDir(sourcePath, targetPath)
		} else if artifacts[0] == "." {
			// copies all source module folder
			_, sourceDir := filepath.Split(sourcePath)
			fullTargetPath := filepath.Join(targetPath, sourceDir)
			dir.CopyDir(sourcePath, fullTargetPath)
		} else {
			//TODO implement other cases of artifacts: subdirectory, file
		}
	} else {
		for _, artifact := range artifacts {
			err := CopyRequiredArtifacts(sourcePath, targetPath, []string{artifact})
			if err != nil {
				return errors.Wrapf(err, "Error copying artifact %v", artifact)
			}
		}
	}
	return nil
}

// validateArtifacts - validates list of required artifacts
func validateArtifacts(ep *dir.MtaLocationParameters, requiredModule *Modules, artifacts []string) error {
	if len(artifacts) == 0 {
		errors.New("No artifacts defined")
	}
	for _, artifact := range artifacts {
		switch true {
		case artifact == "":
			return errors.New("Empty artifact defined")
		case artifact == "." || artifact == "*":
			if len(artifacts) > 1 {
				return errors.New("[*] and [.] artifacts listed among multiple artifacts")
			}
		default:
			//TODO add validations
		}
	}
	return nil
}

// getBuildResultsPath - provides path of build results
func (module *Modules) getBuildResultsPath(ep *dir.MtaLocationParameters) string {
	if module.BuildParams.Path == "" {
		// if no subfolder provided - build results will be saved in the module folder
		return ep.GetSourceModuleDir(module.Path)
	} else {
		// if subfolder provided - build results will be saved in the subfolder of the module folder
		return filepath.Join(ep.GetSourceModuleDir(module.Path), module.BuildParams.Path)
	}
}

// getRequiredTargetPath - provides path of required artifacts
func (requires *BuildRequires) getRequiredTargetPath(ep *dir.MtaLocationParameters, module *Modules) string {
	if requires.TargetPath == "" {
		// if no target folder provided - artifacts will be saved in module folder
		return ep.GetSourceModuleDir(module.Path)
	} else {
		// if target folder provided - artifacts will be saved in the subfolder of the module folder
		return filepath.Join(ep.GetSourceModuleDir(module.Path), requires.TargetPath)
	}
}
