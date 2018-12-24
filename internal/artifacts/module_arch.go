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
func ExecuteBuild(source, target, desc, moduleName, platform string, wdGetter func() (string, error)) error {
	logs.Logger.Infof("build of module <%v> started", moduleName)
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrapf(err, "build of module <%v> failed when initializing location", moduleName)
	}
	err = buildModule(loc, loc, loc.IsDeploymentDescriptor(), moduleName, platform)
	if err != nil {
		return err
	}
	logs.Logger.Infof("build of module <%v> finished successfully", moduleName)
	return nil
}

// ExecutePack - executes packing of module
func ExecutePack(source, target, desc, moduleName, platform string, wdGetter func() (string, error)) error {
	logs.Logger.Infof("pack of module <%v> started", moduleName)

	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrapf(err, "pack of module <%v> failed when initializing location", moduleName)
	}

	module, _, err := commands.GetModuleAndCommands(loc, moduleName)
	if err != nil {
		return errors.Wrapf(err, "pack of module <%v> failed when getting commands", moduleName)
	}

	err = packModule(loc, loc.IsDeploymentDescriptor(), module, moduleName, platform)
	if err != nil {
		return errors.Wrapf(err, "pack of module <%v> failed on module packing", moduleName)
	}

	logs.Logger.Infof("pack of module <%v> finished successfully", moduleName)
	return nil
}

// buildModule - builds module
func buildModule(mtaParser dir.IMtaParser, moduleLoc dir.IModule, deploymentDesc bool, moduleName, platform string) error {

	logs.Logger.Infof("Module %v building started", moduleName)

	// Get module respective command's to execute
	module, mCmd, err := commands.GetModuleAndCommands(mtaParser, moduleName)
	if err != nil {
		return errors.Wrapf(err, "build of module %v failed when getting commands", moduleName)
	}

	if !deploymentDesc {

		// Development descriptor - build includes:
		// 1. module dependencies processing
		e := buildops.ProcessDependencies(mtaParser, moduleLoc, moduleName)
		if e != nil {
			return errors.Wrapf(e, "build of module %v failed when processing dependencies", moduleName)
		}

		// 2. module type dependent commands execution
		modulePath := moduleLoc.GetSourceModuleDir(module.Path)

		// Get module commands
		commands := commands.CmdConverter(modulePath, mCmd)

		// Execute child-process with module respective commands
		e = exec.Execute(commands)
		if e != nil {
			return errors.Wrapf(e, "build of module %v failed when executing commands", moduleName)
		}

		// 3. Packing the modules build artifacts (include node modules)
		// into the artifactsPath dir as data zip
		e = packModule(moduleLoc, false, module, moduleName, platform)
		if e != nil {
			return errors.Wrapf(e, "build of module %v failed when packing", moduleName)
		}
	} else if buildops.PlatformDefined(module, platform) {

		// Deployment descriptor
		// copy module archive to temp directory
		err = copyModuleArchive(moduleLoc, module.Path, moduleName)
		if err != nil {
			return errors.Wrapf(err, "build of module %v failed when copying module's archive", module)
		}
	}
	return nil
}

// packModule - pack build module artifacts
func packModule(ep dir.IModule, deploymentDesc bool, module *mta.Module, moduleName, platform string) error {

	if !buildops.PlatformDefined(module, platform) {
		return nil
	}

	if deploymentDesc {
		return copyModuleArchive(ep, module.Path, moduleName)
	}

	logs.Logger.Infof("pack of module <%v> started", moduleName)
	// Get module relative path
	moduleZipPath := ep.GetTargetModuleDir(moduleName)

	logs.Logger.Info(fmt.Sprintf("module <%v> will be packed and saved in folder %v", moduleName, moduleZipPath))
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	err := os.MkdirAll(moduleZipPath, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "pack of module <%v> failed when creating folder %v", moduleName, moduleZipPath)
	}
	// zipping the build artifacts
	logs.Logger.Infof("zip of module <%v> started", moduleName)
	moduleZipFullPath := moduleZipPath + dataZip
	sourceModuleDir := buildops.GetBuildResultsPath(ep, module)

	err = dir.Archive(sourceModuleDir, moduleZipFullPath)
	if err != nil {
		return errors.Wrapf(err, "pack of module <%v> failed when archiving", moduleName)
	}
	logs.Logger.Infof("pack of module <%v> finished successfully", moduleName)
	return nil
}

// copyModuleArchive - copies module archive to temp directory
func copyModuleArchive(ep dir.IModule, modulePath, moduleName string) error {
	logs.Logger.Infof("copy of archive of module <%v> started", moduleName)
	srcModulePath := ep.GetSourceModuleDir(modulePath)
	moduleSrcZip := filepath.Join(srcModulePath, "data.zip")
	moduleTrgZipPath := ep.GetTargetModuleDir(moduleName)
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	err := os.MkdirAll(moduleTrgZipPath, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "copy of archive of module <%v> failed when creating directory <%v>", moduleName, moduleTrgZipPath)
	}
	moduleTrgZip := filepath.Join(moduleTrgZipPath, "data.zip")
	err = dir.CopyFile(moduleSrcZip, filepath.Join(moduleTrgZipPath, "data.zip"))
	if err != nil {
		return errors.Wrapf(err, "copy of archive of module <%v> failed when copying <%v> to <%v>", moduleName, moduleSrcZip, moduleTrgZip)
	}
	return nil
}
