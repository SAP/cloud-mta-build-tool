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
func ExecuteBuild(source, target, moduleName, platform string, wdGetter func() (string, error)) error {
	logs.Logger.Infof(buildMsg, moduleName)
	loc, err := dir.Location(source, target, dir.Dev, wdGetter)
	if err != nil {
		return errors.Wrapf(err, buildFailedOnLocMsg, moduleName)
	}
	// validate platform
	platform, err = validatePlatform(platform)
	if err != nil {
		return err
	}
	err = buildModule(loc, loc, loc, moduleName, platform)
	if err != nil {
		return err
	}
	return nil
}

// ExecutePack - executes packing of module
func ExecutePack(source, target, moduleName, platform string, wdGetter func() (string, error)) error {
	logs.Logger.Infof(packMsg, moduleName)

	loc, err := dir.Location(source, target, dir.Dev, wdGetter)
	if err != nil {
		return errors.Wrapf(err, packFailedOnLocMsg, moduleName)
	}
	// validate platform
	platform, err = validatePlatform(platform)
	if err != nil {
		return err
	}

	module, _, defaultBuildResult, err := commands.GetModuleAndCommands(loc, moduleName)
	if err != nil {
		return errors.Wrapf(err, packFailedOnCommandsMsg, moduleName)
	}

	err = packModule(loc, loc, module, moduleName, platform, defaultBuildResult)
	if err != nil {
		return err
	}

	return nil
}

// buildModule - builds module
func buildModule(mtaParser dir.IMtaParser, moduleLoc dir.IModule, targetLoc dir.ITargetPath, moduleName, platform string) error {

	// Get module respective command's to execute
	module, mCmd, defaultBuildResults, err := commands.GetModuleAndCommands(mtaParser, moduleName)
	if err != nil {
		return errors.Wrapf(err, buildFailedOnCommandsMsg, moduleName)
	}

	// Development descriptor - build includes:
	// 1. module dependencies processing
	e := buildops.ProcessDependencies(mtaParser, moduleLoc, moduleName)
	if e != nil {
		return errors.Wrapf(e, buildFailedOnDepsMsg, moduleName)
	}

	// 2. module type dependent commands execution
	modulePath := moduleLoc.GetSourceModuleDir(module.Path)

	// Get module commands
	commandList := commands.CmdConverter(modulePath, mCmd)

	// Execute child-process with module respective commands
	e = exec.Execute(commandList)
	if e != nil {
		return errors.Wrapf(e, buildFailedOnExecCmdMsg, moduleName)
	}

	// 3. Packing the modules build artifacts (include node modules)
	// into the artifactsPath dir as data zip
	return packModule(moduleLoc, targetLoc, module, moduleName, platform, defaultBuildResults)
}

// packModule - pack build module artifacts
func packModule(source dir.IModule, target dir.ITargetPath, module *mta.Module, moduleName, platform, defaultBuildResult string) error {

	if !buildops.PlatformDefined(module, platform) {
		return nil
	}

	logs.Logger.Info(fmt.Sprintf(buildResultMsg, moduleName, source.GetTargetModuleDir(moduleName)))

	sourceArtifact, _, _, err := buildops.GetModuleSourceArtifactPath(source, false, module, defaultBuildResult)
	if err != nil {
		return errors.Wrapf(err, packFailedOnBuildArtifactMsg, moduleName)
	}
	targetArtifact, toArchive, err := buildops.GetModuleTargetArtifactPath(source, target, false, module, defaultBuildResult)
	if err != nil {
		return errors.Wrapf(err, packFailedOnTargetArtifactMsg, moduleName)
	}

	if !toArchive {
		return copyModuleArchiveToResultDir(sourceArtifact, targetArtifact, moduleName)
	}

	return archiveModuleToResultDir(sourceArtifact, targetArtifact, getIgnores(module), moduleName)
}

func copyModuleArchiveToResultDir(source, target, moduleName string) error {
	// Create empty folder with name as before the zip process
	// to put the file such as data.zip inside
	modulePathInTmpFolder := filepath.Dir(target)
	err := dir.CreateDirIfNotExist(modulePathInTmpFolder)
	if err != nil {
		return errors.Wrapf(err, packFailedOnFolderCreationMsg, moduleName, modulePathInTmpFolder)
	}

	err = dir.CopyFile(source, target)
	if err != nil {
		return errors.Wrapf(err, packFailedOnCopyMsg, moduleName, source, target)
	}
	return nil
}

func archiveModuleToResultDir(buildResult string, requestedResultFileName string, ignore []string, moduleName string) error {
	// Archive the folder without the ignored files and/or subfolders, which are excluded from the package.
	err := dir.Archive(buildResult, requestedResultFileName, ignore)
	if err != nil {
		return errors.Wrapf(err, PackFailedOnArchMsg, moduleName)
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

// CopyMtaContent copies the content of all modules and resources which are presented in the deployment descriptor,
// in the source directory, to the target directory
func CopyMtaContent(source, target string, copyInParallel bool, wdGetter func() (string, error)) error {

	logs.Logger.Info(copyStartMsg)
	loc, err := dir.Location(source, target, dir.Dep, wdGetter)
	if err != nil {
		return errors.Wrap(err, copyContentFailedOnLocMsg)
	}
	mtaObj, err := loc.ParseFile()
	if err != nil {
		return errors.Wrapf(err, copyContentFailedOnParseMsg, loc.GetMtaYamlPath())
	}
	err = copyModuleContent(loc.GetSource(), loc.GetTargetTmpDir(), mtaObj, copyInParallel)
	if err != nil {
		return err
	}

	err = copyRequiredDependencyContent(loc.GetSource(), loc.GetTargetTmpDir(), mtaObj, copyInParallel)
	if err != nil {
		return err
	}

	return copyResourceContent(loc.GetSource(), loc.GetTargetTmpDir(), mtaObj, copyInParallel)
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
			return handleCopyMtaContentFailure(target, copiedMtaContents, pathNotExistsMsg, []interface{}{mtaPath})
		}
		copiedMtaContents = append(copiedMtaContents, mtaPath)
		targetMtaContent := filepath.Join(target, mtaPath)
		err := copyMtaContentFromPath(sourceMtaContent, targetMtaContent, mtaPath, target, copyInParallel)
		if err != nil {
			return handleCopyMtaContentFailure(target, copiedMtaContents, copyContentFailedMsg, []interface{}{mtaPath, source, err.Error()})
		}
		logs.Logger.Debugf(copyDoneMsg, mtaPath)
	}

	return nil
}

func handleCopyMtaContentFailure(targetLocation string, copiedMtaContents []string,
	message string, messageArguments []interface{}) error {
	errCleanup := cleanUpCopiedContent(targetLocation, copiedMtaContents)
	if errCleanup == nil {
		return errors.Errorf(message, messageArguments...)
	}
	return errors.Errorf(message+cleanupFailedMsg, messageArguments...)
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
	err := dir.CreateDirIfNotExist(filepath.Join(target, mtaContentParentDir))
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
