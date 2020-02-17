package dir

// ModuleLoc - module location type that provides services for stand alone module build command
type ModuleLoc struct {
	loc *Loc
}

// GetTarget - gets the target path;
func (ep *ModuleLoc) GetTarget() string {
	return ep.loc.GetTarget()
}

// GetTargetTmpDir - gets the target path;
func (ep *ModuleLoc) GetTargetTmpDir() string {
	return ep.loc.GetTarget()
}

// GetSourceModuleDir - gets the absolute path to the module
func (ep *ModuleLoc) GetSourceModuleDir(modulePath string) string {
	return ep.loc.GetSourceModuleDir(modulePath)
}

// GetSourceModuleArtifactRelPath - gets the relative path to the module artifact
// The ModuleLoc type is used in context of stand alone module build command and as opposed to the module build command in the context
// of Makefile saves its build result directly under the target (temporary or specific) without considering the original artifact path in the source folder
func (ep *ModuleLoc) GetSourceModuleArtifactRelPath(modulePath, artifactPath string) (string, error) {
	return "", nil
}

// GetTargetModuleDir - gets the to module build results
func (ep *ModuleLoc) GetTargetModuleDir(moduleName string) string {
	return ep.loc.GetTarget()
}

// ModuleLocation - provides target location of stand alone MTA module build result
func ModuleLocation(loc *Loc) *ModuleLoc {
	return &ModuleLoc{loc: loc}
}
