package artifacts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/buildops"
	"cloud-mta-build-tool/internal/commands"
	"cloud-mta-build-tool/internal/exec"
	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/mta"
)

// ExecuteBuild - executes build of module
func ExecuteBuild(source, target, desc, moduleName string, wdGetter func() (string, error)) error {
	logs.Logger.Infof("Build of module  <%v> started", moduleName)
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrapf(err, "Build of module <%v> failed on location initialization", moduleName)
	}
	err = buildModule(loc, loc, loc.IsDeploymentDescriptor(), moduleName)
	if err != nil {
		return errors.Wrapf(err, "Build of module <%v> failed", moduleName)
	}
	logs.Logger.Infof("Build of module  <%v> successfully finished", moduleName)
	return nil
}

// ExecutePack - executes packing of module
func ExecutePack(source, target, desc, moduleName string, wdGetter func() (string, error)) error {
	logs.Logger.Infof("Pack of module  <%v> started", moduleName)

	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrapf(err, "Pack of module <%v> failed on location initialization", moduleName)
	}

	module, _, err := commands.GetModuleAndCommands(loc, moduleName)
	if err != nil {
		return errors.Wrapf(err, "Pack of module <%v> failed on getting modules and commands", moduleName)
	}

	err = packModule(loc, loc.IsDeploymentDescriptor(), module, moduleName)
	if err != nil {
		return errors.Wrapf(err, "Pack of module <%v> failed on module packing", moduleName)
	}

	logs.Logger.Infof("Pack of module  <%v> successfully finished", moduleName)
	return nil
}

// buildModule - builds module
func buildModule(mtaParser dir.IMtaParser, moduleLoc dir.IModule, deploymentDesc bool, moduleName string) error {

	logs.Logger.Infof("Module %v building started", moduleName)

	// Get module respective command's to execute
	module, mCmd, err := commands.GetModuleAndCommands(mtaParser, moduleName)
	if err != nil {
		return errors.Wrapf(err, "Module %v building failed on getting relative path and commands", moduleName)
	}

	if !deploymentDesc {

		// Development descriptor - build includes:
		// 1. module dependencies processing
		e := buildops.ProcessDependencies(mtaParser, moduleLoc, moduleName)
		if e != nil {
			return errors.Wrapf(e, "Module %v building failed on processing dependencies", moduleName)
		}

		// 2. module type dependent commands execution
		modulePath := moduleLoc.GetSourceModuleDir(module.Path)

		// Get module commands
		commands := commands.CmdConverter(modulePath, mCmd)

		// Execute child-process with module respective commands
		e = exec.Execute(commands)
		if e != nil {
			return errors.Wrapf(e, "Module %v building failed on commands execution", moduleName)
		}

		// 3. Packing the modules build artifacts (include node modules)
		// into the artifactsPath dir as data zip
		e = packModule(moduleLoc, false, module, moduleName)
		if e != nil {
			return errors.Wrapf(e, "Module %v building failed on module's packing", moduleName)
		}
	} else if buildops.PlatformsDefined(module) {

		// Deployment descriptor
		// copy module archive to temp directory
		err = copyModuleArchive(moduleLoc, module.Path, moduleName)
		if err != nil {
			return errors.Wrapf(err, "Module %v building failed on module's archive copy", module)
		}
	}
	return nil
}

// packModule - pack build module artifacts
func packModule(ep dir.IModule, deploymentDesc bool, module *mta.Module, moduleName string) error {

	if !buildops.PlatformsDefined(module) {
		return nil
	}

	if deploymentDesc {
		return copyModuleArchive(ep, module.Path, moduleName)
	}

	logs.Logger.Infof("Pack of module %v Started", moduleName)
	// Get module relative path
	moduleZipPath := ep.GetTargetModuleDir(moduleName)

	logs.Logger.Info(fmt.Sprintf("Module %v will be packed and saved in folder %v", moduleName, moduleZipPath))
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	err := os.MkdirAll(moduleZipPath, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "Pack of module %v failed on making directory %v", moduleName, moduleZipPath)
	}
	// zipping the build artifacts
	logs.Logger.Infof("Starting execute zipping module %v ", moduleName)
	moduleZipFullPath := moduleZipPath + dataZip
	sourceModuleDir := buildops.GetBuildResultsPath(ep, module)

	err = dir.Archive(sourceModuleDir, moduleZipFullPath)
	if err != nil {
		return errors.Wrapf(err, "Pack of module %v failed on archiving", moduleName)
	}
	logs.Logger.Infof("Pack of module %v successfully finished", moduleName)
	return nil
}

// copyModuleArchive - copies module archive to temp directory
func copyModuleArchive(ep dir.IModule, modulePath, moduleName string) error {
	logs.Logger.Infof("Copy of module %v archive Started", moduleName)
	srcModulePath := ep.GetSourceModuleDir(modulePath)
	moduleSrcZip := filepath.Join(srcModulePath, "data.zip")
	moduleTrgZipPath := ep.GetTargetModuleDir(moduleName)
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	err := os.MkdirAll(moduleTrgZipPath, os.ModePerm)
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
