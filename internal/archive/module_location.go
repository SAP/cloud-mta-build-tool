package dir

// ModuleLoc - module location type that provides services for stand alone module build command
type ModuleLoc struct {
	loc        *Loc
	targetPath string
}

// GetTarget - gets the target path;
// if it is not provided, use the path of the processed project.
func (ep *ModuleLoc) GetTarget() string {
	return ep.targetPath
}

// GetTargetTmpDir - gets the temporary target directory path.
// The subdirectory in the target folder is named as the source project folder suffixed with "_mta_build_tmp".
// Subdirectory name is prefixed with "." as a hidden folder
// in case of stand alone module build when target folder provided build results will be saved in this target folder
// and not in the temporary folder
func (ep *ModuleLoc) GetTargetTmpDir() string {
	if ep.GetTarget() == "" {
		return ep.loc.GetTargetTmpDir()
	}
	return ep.GetTarget()
}

// GetSourceModuleDir - gets the absolute path to the module
func (ep *ModuleLoc) GetSourceModuleDir(modulePath string) string {
	return ep.loc.GetSourceModuleDir(modulePath)
}

// GetSourceModuleArtifactRelPath - gets the relative path to the module artifact
// The ModuleLoc type is used in context of stand alone module build command and in opposite to the module build command in the context
// of Makefile saves its build result directly under the target (temporary or specific) without considering the original artifact path in the source folder
func (ep *ModuleLoc) GetSourceModuleArtifactRelPath(modulePath, artifactPath string) (string, error) {
	return "", nil
}

// GetTargetModuleDir - gets the to module build results
func (ep *ModuleLoc) GetTargetModuleDir(moduleName string) string {
	if ep.targetPath == "" {
		return ep.loc.GetTargetModuleDir(moduleName)
	}
	return ep.GetTargetTmpDir()
}

// ModuleLocation - provides target location of stand alone MTA module build result
func ModuleLocation(loc *Loc, targetPath string) *ModuleLoc {
	return &ModuleLoc{loc: loc, targetPath: targetPath}
}
