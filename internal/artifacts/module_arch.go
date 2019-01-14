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
	"github.com/SAP/cloud-mta/mta"
)

// ExecuteBuild - executes build of module
func ExecuteBuild(source, target, desc, moduleName, platform string, wdGetter func() (string, error)) error {
	logs.Logger.Infof("build of the %v module started", moduleName)
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrapf(err, "build of the %v module failed when initializing the location", moduleName)
	}
	err = buildModule(loc, loc, loc.IsDeploymentDescriptor(), moduleName, platform)
	if err != nil {
		return err
	}
	logs.Logger.Infof("build of the %v module finished successfully", moduleName)
	return nil
}

// ExecutePack - executes packing of module
func ExecutePack(source, target, desc, moduleName, platform string, wdGetter func() (string, error)) error {
	logs.Logger.Infof("packing of the %v module started", moduleName)

	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrapf(err, "packing of the %v module failed when initializing the location", moduleName)
	}

	module, _, err := commands.GetModuleAndCommands(loc, moduleName)
	if err != nil {
		return errors.Wrapf(err, "packing of the %v module failed when getting commands", moduleName)
	}

	err = packModule(loc, loc.IsDeploymentDescriptor(), module, moduleName, platform)
	if err != nil {
		return err
	}

	logs.Logger.Infof("packing of the %v module finished successfully", moduleName)
	return nil
}

// buildModule - builds module
func buildModule(mtaParser dir.IMtaParser, moduleLoc dir.IModule, deploymentDesc bool, moduleName, platform string) error {

	// Get module respective command's to execute
	module, mCmd, err := commands.GetModuleAndCommands(mtaParser, moduleName)
	if err != nil {
		return errors.Wrapf(err, "build of the %v module failed when getting commands", moduleName)
	}

	if !deploymentDesc {

		// Development descriptor - build includes:
		// 1. module dependencies processing
		e := buildops.ProcessDependencies(mtaParser, moduleLoc, moduleName)
		if e != nil {
			return errors.Wrapf(e, "build of the %v module failed when processing dependencies", moduleName)
		}

		// 2. module type dependent commands execution
		modulePath := moduleLoc.GetSourceModuleDir(module.Path)

		// Get module commands
		commands := commands.CmdConverter(modulePath, mCmd)

		// Execute child-process with module respective commands
		e = exec.Execute(commands)
		if e != nil {
			return errors.Wrapf(e, "build of the %v module failed when executing commands", moduleName)
		}

		// 3. Packing the modules build artifacts (include node modules)
		// into the artifactsPath dir as data zip
		e = packModule(moduleLoc, false, module, moduleName, platform)
		if e != nil {
			return e
		}
	} else if buildops.PlatformDefined(module, platform) {

		// Deployment descriptor
		// copy module archive to temp directory
		err = copyModuleArchive(moduleLoc, module.Path, moduleName)
		if err != nil {
			return errors.Wrapf(err, "build of the %v module failed when copying module's archive", module)
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

	// Get module relative path
	moduleZipPath := ep.GetTargetModuleDir(moduleName)

	logs.Logger.Info(fmt.Sprintf("the %v module will be packed and saved in the %v folder", moduleName, moduleZipPath))
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	err := os.MkdirAll(moduleZipPath, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "packing of the %v module failed when creating the %v folder", moduleName, moduleZipPath)
	}
	// zipping the build artifacts
	logs.Logger.Infof("zip of the %v module started", moduleName)
	moduleZipFullPath := moduleZipPath + dataZip
	sourceModuleDir := buildops.GetBuildResultsPath(ep, module)

	err = dir.Archive(sourceModuleDir, moduleZipFullPath)
	if err != nil {
		return errors.Wrapf(err, "packing of the %v module failed when archiving", moduleName)
	}
	return nil
}

// copyModuleArchive - copies module archive to temp directory
func copyModuleArchive(ep dir.IModule, modulePath, moduleName string) error {
	logs.Logger.Infof("copying of the %v module's archive started", moduleName)
	srcModulePath := ep.GetSourceModuleDir(modulePath)
	moduleSrcZip := filepath.Join(srcModulePath, "data.zip")
	moduleTrgZipPath := ep.GetTargetModuleDir(moduleName)
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	err := os.MkdirAll(moduleTrgZipPath, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "copying of the %v module's archive failed when creating the %v folder", moduleName, moduleTrgZipPath)
	}
	moduleTrgZip := filepath.Join(moduleTrgZipPath, "data.zip")
	err = dir.CopyFile(moduleSrcZip, filepath.Join(moduleTrgZipPath, "data.zip"))
	if err != nil {
		return errors.Wrapf(err, "copying of the %v module's archive failed when copying %v to %v", moduleName, moduleSrcZip, moduleTrgZip)
	}
	logs.Logger.Infof("copying of the %v module's archive finished successfully", moduleName)
	return nil
}

// CopyMtaContent copies the content of all modules and resources which are presented in the deployment descriptor,
// in the source directory, to the target directory
func CopyMtaContent(source, target, desc string, wdGetter func() (string, error)) error {
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrap(err, "copying mta content failed during initialization of deployment descriptor location")
	}
	mta, err := loc.ParseFile()
	if err != nil {
		return errors.Wrap(err, "error while parsing MTA")
	}
	err = copyModuleContent(source, loc.GetTargetTmpDir(), mta)
	if err != nil {
		return err
	}

	err = copyRequiredDependencyContent(source, loc.GetTargetTmpDir(), mta)
	if err != nil {
		return err
	}

	return copyResourceContent(source, loc.GetTargetTmpDir(), mta)
}

func copyModuleContent(source, target string, mta *mta.MTA) error {
	return copyMtaContent(source, target, getModulesWithPaths(mta.Modules))
}

func copyResourceContent(source, target string, mta *mta.MTA) error {
	return copyMtaContent(source, target, getResourcesPaths(mta.Resources))
}

func copyRequiredDependencyContent(source, target string, mta *mta.MTA) error {
	return copyMtaContent(source, target, getRequiredDependencyPaths(mta.Modules))
}

func getRequiredDependencyPaths(mtaModules []*mta.Module) []string {
	result := make([]string, 0)
	for _, module := range mtaModules {
		requiredDependenciesWithPaths := getRequiredDependenciesWithPathsForModule(module)
		result = append(result, requiredDependenciesWithPaths...)
	}
	return result
}

func getRequiredDependenciesWithPathsForModule(module *mta.Module) []string {
	result := make([]string, 0)
	for _, requiredDependency := range module.Requires {
		if requiredDependency.Parameters["path"] != nil {
			result = append(result, requiredDependency.Parameters["path"].(string))
		}
	}
	return result
}
func copyMtaContent(source, target string, mtaPaths []string) error {
	copiendMtaContents := make([]string, 0)
	for _, mtaPath := range mtaPaths {
		mtaContent := filepath.Join(source, mtaPath)
		if doesNotExist(mtaContent) {
			errCleanup := cleanUpCopiedContent(target, copiendMtaContents)
			if errCleanup == nil {
				return fmt.Errorf("%s does not exists in the current location %s", mtaPath, source)
			}
			return fmt.Errorf("%s does not exists in the current location %s; cleanup failed", mtaPath, source)
		}
		copiendMtaContents = append(copiendMtaContents, mtaPath)
		destinationMtaContent := filepath.Join(target, mtaPath)
		err := copyMtaContentFromPath(mtaContent, destinationMtaContent, mtaPath, target)
		if err != nil {
			errCleanup := cleanUpCopiedContent(target, copiendMtaContents)
			if errCleanup == nil {
				return fmt.Errorf("Error copying mta content %s to target directory %s: %s", mtaContent, destinationMtaContent, err.Error())
			}
			return fmt.Errorf("Error copying mta content %s to target directory %s: %s; cleanup failed", mtaContent, destinationMtaContent, err.Error())
		}
	}

	return nil
}

func copyMtaContentFromPath(mtaContent, destinationMtaContent, mtaContentPath, target string) error {
	mtaContentInfo, _ := os.Stat(mtaContent)
	if mtaContentInfo.IsDir() {
		err := os.MkdirAll(destinationMtaContent, os.ModePerm)
		if err != nil {
			return err
		}
		return dir.CopyDir(mtaContent, destinationMtaContent)
	}

	mtaContentParentDir := filepath.Dir(mtaContentPath)
	err := os.MkdirAll(filepath.Join(target, mtaContentParentDir), os.ModePerm)
	if err != nil {
		return err
	}
	return dir.CopyFile(mtaContent, destinationMtaContent)
}

func cleanUpCopiedContent(targetLocation string, copiendMtaContents []string) error {
	for _, copiedMtaContent := range copiendMtaContents {
		err := os.RemoveAll(filepath.Join(targetLocation, copiedMtaContent))
		if err != nil {
			return err
		}
	}
	return nil
}

func doesNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func getModulesWithPaths(mtaModules []*mta.Module) []string {
	result := make([]string, 0)
	for _, module := range mtaModules {
		if module.Path != "" {
			result = append(result, module.Path)
		}
	}
	return result
}

func getResourcesPaths(resources []*mta.Resource) []string {
	result := make([]string, 0)
	for _, resource := range resources {
		if resource.Parameters["path"] != nil {
			result = append(result, resource.Parameters["path"].(string))
		}
	}
	return result
}
