package mta

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fsys"
)

// Order of modules building is done according to the dependencies defined in build parameters.
// In case of problems in this definition build process should not start and corresponding error must be provided.
// Possible problems:
// 1.	Cyclic dependencies
// 2.	Dependency on not defined module
// 3.	Wrong definition of artifacts

// ProcessRequirements - Processes build requirement of module (using moduleName).
func (requires *BuildRequires) ProcessRequirements(ep *dir.MtaLocationParameters, mta *MTA, moduleName string) error {
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
	artifactsStr := strings.Replace(requires.Artifacts, "[", "", 1)
	artifactsStr = strings.Replace(artifactsStr, "]", "", 1)
	artifacts := strings.Split(artifactsStr, ",")

	// Validate artifacts
	err = validateArtifacts(ep, requiredModule, artifacts)
	if err != nil {
		return errors.Wrapf(err, "Processing requirements of module %v based on module %v failed on Artifacts validation", moduleName, requiredModule.Name)
	}
	// Build paths for artifacts copying
	sourcePath, err := requiredModule.getBuildResultsPath(ep)
	if err != nil {
		return errors.Wrapf(err, "Processing requirements of module %v based on module %v failed on getting Results Path", moduleName, requiredModule.Name)
	}
	targetPath, err := requires.getRequiredTargetPath(ep, module)
	if err != nil {
		return errors.Wrapf(err, "Processing requirements of module %v based on module %v failed on getting Required Target Path", moduleName, requiredModule.Name)
	}
	// execute copy of artifacts
	err = copyRequiredArtifacts(sourcePath, targetPath, artifacts)

	if err != nil {
		return errors.Wrapf(err, "Processing requirements of module %v based on module %v failed on artifacts copying", moduleName, requiredModule.Name)
	}
	return nil
}

// copyRequiredArtifacts - copies artifacts of predecessor (source module) to dependent (target module)
func copyRequiredArtifacts(sourcePath, targetPath string, artifacts []string) error {
	if len(artifacts) == 1 {
		if artifacts[0] == "*" {
			// copies all source module folder's entries
			if err := dir.CopyDir(sourcePath, targetPath); err != nil {
				return errors.Wrapf(err, "Error copying dir")
			}
		} else if artifacts[0] == "." {
			// copies all source module folder
			_, sourceDir := filepath.Split(sourcePath)
			fullTargetPath := filepath.Join(targetPath, sourceDir)
			err := dir.CopyDir(sourcePath, fullTargetPath)
			if err != nil {
				return errors.Wrapf(err, "Error copying dir")
			}
		} else {
			// TODO implement other cases of artifacts: subdirectory, file
		}
	} else {
		for _, artifact := range artifacts {
			err := copyRequiredArtifacts(sourcePath, targetPath, []string{artifact})
			if err != nil {
				return errors.Wrapf(err, "Error copying artifact %v", artifact)
			}
		}
	}
	return nil
}

// validateArtifacts - validates list of required artifacts
//noinspection GoUnusedParameter,GoUnusedParameter
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
			// TODO add validations
		}
	}
	return nil
}

// getBuildResultsPath - provides path of build results
func (module *Modules) getBuildResultsPath(ep *dir.MtaLocationParameters) (string, error) {
	if module.BuildParams.Path == "" {
		// if no sub-folder provided - build results will be saved in the module folder
		return ep.GetSourceModuleDir(module.Path)
	}
	// if sub-folder provided - build results will be saved in the subfolder of the module folder
	source, err := ep.GetSourceModuleDir(module.Path)
	if err != nil {
		return "", errors.Wrap(err, "getBuildResultsPath failed")
	}
	return filepath.Join(source, module.BuildParams.Path), nil
}

// getRequiredTargetPath - provides path of required artifacts
func (requires *BuildRequires) getRequiredTargetPath(ep *dir.MtaLocationParameters, module *Modules) (string, error) {
	if requires.TargetPath == "" {
		// if no target folder provided - artifacts will be saved in module folder
		return ep.GetSourceModuleDir(module.Path)
	}
	// if target folder provided - artifacts will be saved in the sub-folder of the module folder
	source, err := ep.GetSourceModuleDir(module.Path)
	if err != nil {
		return "", errors.Wrap(err, "getRequiredTargetPath failed")
	}
	return filepath.Join(source, requires.TargetPath), nil
}
