package buildops

import (
	"path/filepath"

	"cloud-mta-build-tool/mta"
	"github.com/pkg/errors"

	fs "cloud-mta-build-tool/internal/fsys"
)

// Order of modules building is done according to the dependencies defined in build parameters.
// In case of problems in this definition build process should not start and corresponding error must be provided.
// Possible problems:
// 1.	Cyclic dependencies
// 2.	Dependency on not defined module

// ProcessRequirements - Processes build requirement of module (using moduleName).
func ProcessRequirements(ep *fs.Loc, mta *mta.MTA, requires *mta.BuildRequires, moduleName string) error {
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
	err = fs.CopyByPatterns(sourcePath, targetPath, artifacts)

	if err != nil {
		return errors.Wrapf(err, "Processing requirements of module %v based on module %v failed on artifacts copying", moduleName, requiredModule.Name)
	}
	return nil
}

// getBuildResultsPath - provides path of build results
func getBuildResultsPath(ep *fs.Loc, module *mta.Module) (string, error) {
	path, err := ep.GetSourceModuleDir(module.Path)
	if err != nil {
		return "", errors.Wrapf(err, "getBuildResultsPath failed getting directory of module %v", module.Path)
	}
	// if no sub-folder provided - build results will be saved in the module folder
	if module.BuildParams.Path != "" {
		// if sub-folder provided - build results are located in the subfolder of the module folder
		path = filepath.Join(path, module.BuildParams.Path)
	}
	return path, nil
}

// getRequiredTargetPath - provides path of required artifacts
func getRequiredTargetPath(ep *fs.Loc, module *mta.Module, requires *mta.BuildRequires) (string, error) {
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
