package artifacts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/builders"
	"cloud-mta-build-tool/internal/buildops"
	"cloud-mta-build-tool/internal/exec"
	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/mta"
)

// BuildModule - builds module
func BuildModule(ep *dir.Loc, moduleName string) error {

	logs.Logger.Infof("Module %v building started", moduleName)

	// Get module respective command's to execute
	module, mCmd, err := builders.GetModuleAndCommands(ep, moduleName)
	if err != nil {
		return errors.Wrapf(err, "Module %v building failed on getting relative path and commands", moduleName)
	}

	if !ep.IsDeploymentDescriptor() {

		// Development descriptor - build includes:
		// 1. module dependencies processing
		e := buildops.ProcessDependencies(ep, moduleName)
		if e != nil {
			return errors.Wrapf(e, "Module %v building failed on processing dependencies", moduleName)
		}

		// 2. module type dependent commands execution
		modulePath, e := ep.GetSourceModuleDir(module.Path)
		if e != nil {
			return errors.Wrapf(e, "Module %v building failed on getting source module directory", moduleName)
		}

		// Get module commands
		commands := builders.CmdConverter(modulePath, mCmd)

		// Execute child-process with module respective commands
		e = exec.Execute(commands)
		if e != nil {
			return errors.Wrapf(e, "Module %v building failed on commands execution", moduleName)
		}

		// 3. Packing the modules build artifacts (include node modules)
		// into the artifactsPath dir as data zip
		e = PackModule(ep, module, moduleName)
		if e != nil {
			return errors.Wrapf(e, "Module %v building failed on module's packing", moduleName)
		}
	} else if buildops.PlatformsDefined(module) {

		// Deployment descriptor
		// copy module archive to temp directory
		err = CopyModuleArchive(ep, module.Path, moduleName)
		if err != nil {
			return errors.Wrapf(err, "Module %v building failed on module's archive copy", module)
		}
	}
	return nil
}

// PackModule - pack build module artifacts
func PackModule(ep *dir.Loc, module *mta.Module, moduleName string) error {

	if !buildops.PlatformsDefined(module) {
		return nil
	}

	if ep.IsDeploymentDescriptor() {
		return CopyModuleArchive(ep, module.Path, moduleName)
	}

	logs.Logger.Infof("Pack of module %v Started", moduleName)
	// Get module relative path
	moduleZipPath, err := ep.GetTargetModuleDir(moduleName)
	if err != nil {
		return errors.Wrapf(err, "Pack of module %v failed getting target module directory", moduleName)
	}
	logs.Logger.Info(fmt.Sprintf("Module %v will be packed and saved in folder %v", moduleName, moduleZipPath))
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	err = os.MkdirAll(moduleZipPath, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "Pack of module %v failed on making directory %v", moduleName, moduleZipPath)
	}
	// zipping the build artifacts
	logs.Logger.Infof("Starting execute zipping module %v ", moduleName)
	moduleZipFullPath := moduleZipPath + dataZip
	sourceModuleDir, err := buildops.GetBuildResultsPath(ep, module)
	if err != nil {
		return errors.Wrapf(err, "Pack of module %v failed on getting source module directory with relative path %v",
			moduleName, module.Path)
	}
	if err = dir.Archive(sourceModuleDir, moduleZipFullPath); err != nil {
		return errors.Wrapf(err, "Pack of module %v failed on archiving", moduleName)
	}
	logs.Logger.Infof("Pack of module %v successfully finished", moduleName)
	return nil
}

// CopyModuleArchive - copies module archive to temp directory
func CopyModuleArchive(ep *dir.Loc, modulePath, moduleName string) error {
	logs.Logger.Infof("Copy of module %v archive Started", moduleName)
	srcModulePath, err := ep.GetSourceModuleDir(modulePath)
	if err != nil {
		return errors.Wrapf(err, "Copy of module %v archive failed getting source module directory", moduleName)
	}
	moduleSrcZip := filepath.Join(srcModulePath, "data.zip")
	moduleTrgZipPath, err := ep.GetTargetModuleDir(moduleName)
	if err != nil {
		return errors.Wrapf(err, "Copy of module %v archive failed getting target module directory", moduleName)
	}
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	err = os.MkdirAll(moduleTrgZipPath, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "Copy of module %v archive on making directory %v", moduleName, moduleTrgZipPath)
	}
	moduleTrgZip := filepath.Join(moduleTrgZipPath, "data.zip")
	err = dir.CopyFile(moduleSrcZip, filepath.Join(moduleTrgZipPath, "data.zip"))
	if err != nil {
		return errors.Wrapf(err, "Copy of module %v archive failed copying %v to %v", moduleName, moduleSrcZip, moduleTrgZip)
	}
	return nil
}
