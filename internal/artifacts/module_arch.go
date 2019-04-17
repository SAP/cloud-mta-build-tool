package artifacts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta/mta"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/buildops"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/exec"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

const (
	ignore = "ignore"
)

// ExecuteBuild - executes build of module
func ExecuteBuild(source, target, desc, moduleName, platform string, wdGetter func() (string, error)) error {
	logs.Logger.Infof(`building the "%v" module...`, moduleName)
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrapf(err, `build of the "%v" module failed when initializing the location`, moduleName)
	}
	// validate platform
	platform, err = validatePlatform(platform)
	if err != nil {
		return err
	}
	err = buildModule(loc, loc, loc.IsDeploymentDescriptor(), moduleName, platform)
	if err != nil {
		return err
	}
	return nil
}

// ExecutePack - executes packing of module
func ExecutePack(source, target, desc, moduleName, platform string, wdGetter func() (string, error)) error {
	logs.Logger.Infof(`packing the "%v" module...`, moduleName)

	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrapf(err, `packing of the "%v" module failed when initializing the location`, moduleName)
	}
	// validate platform
	platform, err = validatePlatform(platform)
	if err != nil {
		return err
	}

	module, _, defaultBuildResult, err := commands.GetModuleAndCommands(loc, moduleName)
	if err != nil {
		return errors.Wrapf(err, `packing of the "%v" module failed when getting commands`, moduleName)
	}

	err = packModule(loc, loc.IsDeploymentDescriptor(), module, moduleName, platform, defaultBuildResult)
	if err != nil {
		return err
	}

	return nil
}

// buildModule - builds module
func buildModule(mtaParser dir.IMtaParser, moduleLoc dir.IModule, deploymentDesc bool, moduleName, platform string) error {

	// Get module respective command's to execute
	module, mCmd, defaultBuildResults, err := commands.GetModuleAndCommands(mtaParser, moduleName)
	if err != nil {
		return errors.Wrapf(err, `build of the "%v" module failed when getting commands`, moduleName)
	}

	if !deploymentDesc {

		// Development descriptor - build includes:
		// 1. module dependencies processing
		e := buildops.ProcessDependencies(mtaParser, moduleLoc, moduleName)
		if e != nil {
			return errors.Wrapf(e, `build of the "%v" module failed when processing dependencies`, moduleName)
		}

		// 2. module type dependent commands execution
		modulePath := moduleLoc.GetSourceModuleDir(module.Path)

		// Get module commands
		commands := commands.CmdConverter(modulePath, mCmd)

		// Execute child-process with module respective commands
		e = exec.Execute(commands)
		if e != nil {
			return errors.Wrapf(e, `build of the "%v" module failed when executing commands`, moduleName)
		}

		// 3. Packing the modules build artifacts (include node modules)
		// into the artifactsPath dir as data zip
		e = packModule(moduleLoc, false, module, moduleName, platform, defaultBuildResults)
		if e != nil {
			return e
		}
	} else if buildops.PlatformDefined(module, platform) {

		// Deployment descriptor
		// copy module archive to temp directory
		err = copyModuleArchive(moduleLoc, module.Path, moduleName)
		if err != nil {
			return errors.Wrapf(err, `build of the "%v" module failed when copying module's archive`, module)
		}
	}
	return nil
}

// packModule - pack build module artifacts
func packModule(ep dir.IModule, deploymentDesc bool, module *mta.Module, moduleName, platform, defaultBuildResult string) error {

	if !buildops.PlatformDefined(module, platform) {
		return nil
	}

	if deploymentDesc {
		return copyModuleArchive(ep, module.Path, moduleName)
	}
	// Get module relative path
	moduleZipPath := ep.GetTargetModuleDir(moduleName)

	logs.Logger.Info(fmt.Sprintf(`the build results of the "%v" module will be packed and saved in the "%v" folder`, moduleName, moduleZipPath))
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	err := os.MkdirAll(moduleZipPath, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, `packing of the "%v" module failed when creating the "%v" folder`, moduleName, moduleZipPath)
	}
	// zipping the build artifacts
	logs.Logger.Infof(`zipping the %v module...`, moduleName)
	buildResults, _, err := buildops.GetBuildResultsPath(ep, module, defaultBuildResult)
	if err != nil {
		return errors.Wrapf(err, `packing the "%v" module failed while getting the build results path`,
			moduleName)
	}

	entry, err := os.Stat(buildResults)
	if err != nil {
		return errors.Wrapf(err, `packing the "%v" module failed; the "%v" build results path does not exist`,
			moduleName, buildResults)
	}
	if !entry.IsDir() && isArchive(buildResults) {
		err = dir.CopyFile(buildResults, filepath.Join(moduleZipPath, filepath.Base(buildResults)))
		if err != nil {
			return errors.Wrapf(err, `packing of the "%v" module failed when copying the "%s" path to the "%s" folder`,
				moduleName, buildResults, moduleZipPath)
		}
		return nil
	}

	moduleZipFullPath := moduleZipPath + dataZip
	// get ignore - get files and/or subfolders to exclude from the package.
	ignore := getIgnores(module)
	err = dir.Archive(buildResults, moduleZipFullPath, ignore)
	if err != nil {
		return errors.Wrapf(err, `packing of the "%v" module failed when archiving`, moduleName)
	}
	return nil
}

// getIgnores - get files and/or subfolders to exclude from the package.
func getIgnores(module *mta.Module) []string {
	var ignoreList []string
	// ignore defined in build params is declared
	if module.BuildParams != nil && module.BuildParams[ignore] != nil {
		ignoreList = convert(module.BuildParams[ignore].([]interface{}))
	}

	return ignoreList
}

// Convert slice []interface{} to slice []string
func convert(data []interface{}) []string {
	aString := make([]string, len(data))
	for i, v := range data {
		aString[i] = v.(string)
	}
	return aString
}

func isArchive(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".zip" || ext == ".jar" || ext == ".war"
}

// copyModuleArchive - copies module archive to temp directory
func copyModuleArchive(ep dir.IModule, modulePath, moduleName string) error {
	logs.Logger.Infof(`copying the "%v" module's archive`, moduleName)
	srcModulePath := ep.GetSourceModuleDir(modulePath)
	moduleSrcZip := filepath.Join(srcModulePath, "data.zip")
	moduleTrgZipPath := ep.GetTargetModuleDir(moduleName)
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	err := os.MkdirAll(moduleTrgZipPath, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, `copying of the "%v" module's archive failed when creating the "%v" folder`, moduleName, moduleTrgZipPath)
	}
	moduleTrgZip := filepath.Join(moduleTrgZipPath, "data.zip")
	err = dir.CopyFile(moduleSrcZip, filepath.Join(moduleTrgZipPath, "data.zip"))
	if err != nil {
		return errors.Wrapf(err, `copying of the "%v" module's archive failed when copying "%v" to "%v"`, moduleName, moduleSrcZip, moduleTrgZip)
	}
	return nil
}

// CopyMtaContent copies the content of all modules and resources which are presented in the deployment descriptor,
// in the source directory, to the target directory
func CopyMtaContent(source, target, desc string, copyInParallel bool, wdGetter func() (string, error)) error {

	logs.Logger.Info("copying the MTA content...")
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrap(err,
			"copying the MTA content failed during the initialization of deployment descriptor location")
	}
	mta, err := loc.ParseFile()
	if err != nil {
		return errors.Wrapf(err, `copying the MTA content failed when parsing the %s file`, loc.GetMtaYamlPath())
	}
	err = copyModuleContent(loc.GetSource(), loc.GetTargetTmpDir(), mta, copyInParallel)
	if err != nil {
		return err
	}

	err = copyRequiredDependencyContent(loc.GetSource(), loc.GetTargetTmpDir(), mta, copyInParallel)
	if err != nil {
		return err
	}

	return copyResourceContent(loc.GetSource(), loc.GetTargetTmpDir(), mta, copyInParallel)
}

func copyModuleContent(source, target string, mta *mta.MTA, copyInParallel bool) error {
	return copyMtaContent(source, target, getModulesWithPaths(mta.Modules), copyInParallel)
}

func copyResourceContent(source, target string, mta *mta.MTA, copyInParallel bool) error {
	return copyMtaContent(source, target, getResourcesPaths(mta.Resources), copyInParallel)
}

func copyRequiredDependencyContent(source, target string, mta *mta.MTA, copyInParallel bool) error {
	return copyMtaContent(source, target, getRequiredDependencyPaths(mta.Modules), copyInParallel)
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
func copyMtaContent(source, target string, mtaPaths []string, copyInParallel bool) error {
	copiedMtaContents := make([]string, 0)
	for _, mtaPath := range mtaPaths {
		sourceMtaContent := filepath.Join(source, mtaPath)
		if doesNotExist(sourceMtaContent) {
			return handleCopyMtaContentFailure(target, copiedMtaContents,
				`"%s" does not exist in the MTA project location`, []interface{}{mtaPath})
		}
		copiedMtaContents = append(copiedMtaContents, mtaPath)
		targetMtaContent := filepath.Join(target, mtaPath)
		err := copyMtaContentFromPath(sourceMtaContent, targetMtaContent, mtaPath, target, copyInParallel)
		if err != nil {
			return handleCopyMtaContentFailure(target, copiedMtaContents,
				`error copying the "%s" MTA content to the "%s" target directory because: %s`, []interface{}{mtaPath, source, err.Error()})
		}
		logs.Logger.Debugf(`copied "%s"`, mtaPath)
	}

	return nil
}

func handleCopyMtaContentFailure(targetLocation string, copiedMtaContents []string,
	message string, messageArguments []interface{}) error {
	errCleanup := cleanUpCopiedContent(targetLocation, copiedMtaContents)
	if errCleanup == nil {
		return fmt.Errorf(message, messageArguments...)
	}
	return fmt.Errorf(message+"; cleanup failed", messageArguments...)
}

func copyMtaContentFromPath(sourceMtaContent, targetMtaContent, mtaContentPath, target string, copyInParallel bool) error {
	mtaContentInfo, _ := os.Stat(sourceMtaContent)
	if mtaContentInfo.IsDir() {
		if copyInParallel {
			return dir.CopyDir(sourceMtaContent, targetMtaContent, true, dir.CopyEntriesInParallel)
		}
		return dir.CopyDir(sourceMtaContent, targetMtaContent, true, dir.CopyEntries)
	}

	mtaContentParentDir := filepath.Dir(mtaContentPath)
	err := os.MkdirAll(filepath.Join(target, mtaContentParentDir), os.ModePerm)
	if err != nil {
		return err
	}
	return dir.CopyFileWithMode(sourceMtaContent, targetMtaContent, mtaContentInfo.Mode())
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
