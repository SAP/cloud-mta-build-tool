package artifacts

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/buildops"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/exec"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"github.com/SAP/cloud-mta/mta"
)

const (
	ignore = "ignore"
)

// ExecuteBuild - executes build of module from Makefile
func ExecuteBuild(source, target string, extensions []string, moduleName, platform string, wdGetter func() (string, error)) error {
	if moduleName == "" {
		return errors.New(buildFailedOnEmptyModuleMsg)
	}

	logs.Logger.Infof(buildMsg, moduleName)
	loc, err := dir.Location(source, target, dir.Dev, extensions, wdGetter)
	if err != nil {
		return errors.Wrapf(err, buildFailedMsg, moduleName)
	}
	err = buildModule(loc, loc, moduleName, platform, true, true, map[string]string{})
	if err != nil {
		return err
	}
	logs.Logger.Infof(buildFinishedMsg, moduleName)
	return nil
}

// ExecuteSoloBuild - executes build of module from stand alone command
func ExecuteSoloBuild(source, target string, extensions []string, modulesNames []string, allDependencies bool,
	generateMtadFlag bool, platform string,
	wdGetter func() (string, error)) error {

	if len(modulesNames) == 0 {
		return errors.New(buildFailedOnEmptyModulesMsg)
	}

	sourceDir, err := getSoloModuleBuildAbsSource(source, wdGetter)
	if err != nil {
		return wrapBuildError(err, modulesNames)
	}

	loc, err := dir.Location(sourceDir, "", dir.Dev, extensions, wdGetter)
	if err != nil {
		return wrapBuildError(err, modulesNames)
	}

	mtaObj, err := loc.ParseFile()
	if err != nil {
		return err
	}

	allModulesSorted, err := buildops.GetModulesNames(mtaObj)
	if err != nil {
		return wrapBuildError(err, modulesNames)
	}

	selectedModulesMap := make(map[string]bool)
	var selectedModulesWithDependenciesMap map[string]bool

	for _, moduleName := range modulesNames {
		selectedModulesMap[moduleName] = true
	}

	if allDependencies {
		selectedModulesWithDependenciesMap = make(map[string]bool)
		for module := range selectedModulesMap {
			err = collectSelectedModulesAndDependencies(mtaObj, selectedModulesWithDependenciesMap, module)
			if err != nil {
				return wrapBuildError(err, modulesNames)
			}
		}
	} else {
		selectedModulesWithDependenciesMap = selectedModulesMap
	}

	sortedModules := sortModules(allModulesSorted, selectedModulesWithDependenciesMap)

	if allDependencies && len(sortedModules) > 1 {
		logs.Logger.Infof(buildWithDependenciesMsg, `"`+strings.Join(sortedModules, `","`)+`"`)
	} else if len(sortedModules) > 1 {
		logs.Logger.Infof(multiBuildMsg, `"`+strings.Join(sortedModules, `", "`)+`"`)
	}

	packedModulePaths, err := buildModules(sourceDir, target, extensions, sortedModules, selectedModulesMap, wdGetter)
	if err != nil {
		return wrapBuildError(err, modulesNames)
	}

	if generateMtadFlag {
		err = generateMtad(mtaObj, loc, target, platform, packedModulePaths, wdGetter)
		if err != nil {
			return wrapBuildError(err, modulesNames)
		}
	}

	if len(modulesNames) > 1 {
		logs.Logger.Infof(multiBuildFinishedMsg)
	}

	return nil
}

func generateMtad(mtaObj *mta.MTA, loc dir.ITargetPath, target string,
	platform string, packedModulePaths map[string]string, wdGetter func() (string, error)) error {

	platform, err := validatePlatform(platform)
	if err != nil {
		return err
	}

	mtadTargetPath, err := getMtadPath(target, wdGetter)
	if err != nil {
		return err
	}
	mtadLocation := mtadLoc{path: mtadTargetPath}

	return genMtad(mtaObj, &mtadLocation, loc, false, platform, false, packedModulePaths, yaml.Marshal)
}

func getMtadPath(target string, wdGetter func() (string, error)) (string, error) {
	if target != "" {
		return target, nil
	}
	return wdGetter()
}

func wrapBuildError(err error, modules []string) error {
	if len(modules) == 1 {
		return errors.Wrapf(err, buildFailedMsg, modules[0])
	}
	return errors.Wrapf(err, multiBuildFailedMsg)
}

func collectSelectedModulesAndDependencies(mtaObj *mta.MTA, modulesWithDependencies map[string]bool, moduleName string) error {

	if modulesWithDependencies[moduleName] {
		return nil
	}

	modulesWithDependencies[moduleName] = true
	module, err := mtaObj.GetModuleByName(moduleName)
	if err != nil {
		return err
	}
	for _, requires := range buildops.GetBuildRequires(module) {
		requiredModule, err := mtaObj.GetModuleByName(requires.Name)
		if err != nil {
			return err
		}

		err = collectSelectedModulesAndDependencies(mtaObj, modulesWithDependencies, requiredModule.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildModules(source, target string, extensions []string, modulesToBuild []string,
	modulesToPack map[string]bool, wdGetter func() (string, error)) (packedModulePaths map[string]string, err error) {

	buildResults := make(map[string]string)
	for _, module := range modulesToBuild {
		err := buildSelectedModule(source, target, extensions, module, modulesToPack[module], buildResults, wdGetter)

		if err != nil {
			return nil, err
		}
	}

	packedModulePaths = make(map[string]string)
	for buildResult, moduleName := range buildResults {
		packedModulePaths[moduleName] = buildResult
	}
	return packedModulePaths, nil
}

func buildSelectedModule(source, target string, extensions []string, module string,
	toPack bool, buildResults map[string]string, wdGetter func() (string, error)) error {

	logs.Logger.Infof(buildMsg, module)

	moduleLoc, err := getModuleLocation(source, target, module, extensions, wdGetter)
	if err != nil {
		return err
	}

	err = buildModule(moduleLoc, moduleLoc, module, "", false, toPack, buildResults)
	if err != nil {
		return err
	}

	logs.Logger.Infof(buildFinishedMsg, module)
	return nil
}

func sortModules(allModulesSorted []string, selectedModulesMap map[string]bool) []string {
	var result []string
	for _, module := range allModulesSorted {
		_, selected := selectedModulesMap[module]
		if selected {
			result = append(result, module)
		}
	}
	return result
}

func getModuleLocation(source, target, moduleName string, extensions []string, wdGetter func() (string, error)) (*dir.ModuleLoc, error) {
	targetDir, err := getSoloModuleBuildAbsTarget(source, target, moduleName, wdGetter)
	if err != nil {
		return nil, err
	}

	loc, err := dir.Location(source, targetDir, dir.Dev, extensions, wdGetter)
	if err != nil {
		return nil, err
	}

	return dir.ModuleLocation(loc), nil
}

func getSoloModuleBuildAbsSource(source string, wdGetter func() (string, error)) (string, error) {
	if source == "" {
		return wdGetter()
	}
	return filepath.Abs(source)
}

func getSoloModuleBuildAbsTarget(absSource, target, moduleName string, wdGetter func() (string, error)) (string, error) {
	if target != "" {
		return filepath.Abs(target)
	}

	target, err := wdGetter()
	if err != nil {
		return "", err
	}
	_, projectFolderName := filepath.Split(absSource)
	tmpFolderName := "." + projectFolderName + dir.TempFolderSuffix

	// default target is <current folder>/.<project folder>_mta_tmp/<module_name>
	return filepath.Join(target, tmpFolderName, moduleName), nil
}

// ExecutePack - executes packing of module
func ExecutePack(source, target string, extensions []string, moduleName, platform string, wdGetter func() (string, error)) error {
	logs.Logger.Infof(packMsg, moduleName)

	loc, err := dir.Location(source, target, dir.Dev, extensions, wdGetter)
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

	if buildops.IfNoSource(module) {
		logs.Logger.Infof(packSkippedMsg, module.Name)
		return nil
	}

	if module.Path == "" {
		return fmt.Errorf(packFailedOnEmptyPathMsg, moduleName)
	}

	err = packModule(loc, module, moduleName, platform, defaultBuildResult, true, map[string]string{})
	if err != nil {
		return err
	}

	return nil
}

// buildModule - builds module
func buildModule(mtaParser dir.IMtaParser, moduleLoc dir.IModule, moduleName, platform string,
	checkPlatform bool, toPack bool, buildResults map[string]string) error {

	var err error
	if checkPlatform {
		// validate platform
		platform, err = validatePlatform(platform)
		if err != nil {
			return err
		}
	}

	// Get module respective command's to execute
	module, mCmd, defaultBuildResults, err := commands.GetModuleAndCommands(mtaParser, moduleName)
	if err != nil {
		return errors.Wrapf(err, buildFailedOnCommandsMsg, moduleName)
	}

	if buildops.IfNoSource(module) {
		logs.Logger.Infof(buildSkippedMsg, module.Name)
		return nil
	}

	if module.Path == "" {
		return fmt.Errorf(buildFailedOnEmptyPathMsg, moduleName)
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
	commandList, e := commands.CmdConverter(modulePath, mCmd)
	if e != nil {
		return errors.Wrapf(e, buildFailedOnCommandsMsg, moduleName)
	}

	// Execute child-process with module respective commands
	var timeout string
	if module.BuildParams != nil && module.BuildParams["timeout"] != nil {
		var ok bool
		timeout, ok = module.BuildParams["timeout"].(string)
		if !ok {
			return errors.Errorf(exec.ExecInvalidTimeoutMsg, fmt.Sprint(module.BuildParams["timeout"]))
		}
	}
	e = exec.ExecuteWithTimeout(commandList, timeout, true)
	if e != nil {
		return errors.Wrapf(e, buildFailedMsg, moduleName)
	}

	if toPack {
		// 3. Packing the modules build artifacts (include node modules)
		// into the artifactsPath dir as data zip
		return packModule(moduleLoc, module, moduleName, platform, defaultBuildResults, checkPlatform, buildResults)
	}

	return nil
}

// packModule - pack build module artifacts
func packModule(moduleLoc dir.IModule, module *mta.Module, moduleName, platform,
defaultBuildResult string, checkPlatform bool, buildResults map[string]string) error {

	if checkPlatform && !buildops.PlatformDefined(module, platform) {
		return nil
	}

	logs.Logger.Info(fmt.Sprintf(buildResultMsg, moduleName, moduleLoc.GetTargetModuleDir(moduleName)))

	sourceArtifact, err := buildops.GetModuleSourceArtifactPath(moduleLoc, false, module, defaultBuildResult, true)
	if err != nil {
		return errors.Wrapf(err, packFailedOnBuildArtifactMsg, moduleName)
	}
	targetArtifact, toArchive, err := buildops.GetModuleTargetArtifactPath(moduleLoc, false, module, defaultBuildResult)
	if err != nil {
		return errors.Wrapf(err, packFailedOnTargetArtifactMsg, moduleName)
	}

	conflictingModule, ok := buildResults[targetArtifact]
	if ok {
		return fmt.Errorf(multiBuildWithPathsConflictMsg, conflictingModule, module.Name, filepath.Dir(targetArtifact), filepath.Base(targetArtifact))
	}
	buildResults[targetArtifact] = moduleName

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
func CopyMtaContent(source, target string, extensions []string, copyInParallel bool, wdGetter func() (string, error)) error {

	logs.Logger.Info(copyStartMsg)
	loc, err := dir.Location(source, target, dir.Dep, extensions, wdGetter)
	if err != nil {
		return errors.Wrap(err, copyContentFailedOnLocMsg)
	}
	mtaObj, err := loc.ParseFile()
	if err != nil {
		return errors.Wrapf(err, copyContentFailedMsg)
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
			return handleCopyMtaContentFailure(target, copiedMtaContents, copyContentCopyFailedMsg, []interface{}{mtaPath, source, err.Error()})
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
	return errors.Errorf(message+"; "+cleanupFailedMsg, messageArguments...)
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
