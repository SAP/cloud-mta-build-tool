package buildops

const (
	// WrongBuildResultMsg - message raised on wrong build result
	WrongBuildResultMsg = `the build result must be a string; change "%v" in the "%s" module for a string value`
	// WrongBuildArtifactNameMsg - message raised on wrong build artifact name
	WrongBuildArtifactNameMsg = `the build artifact name must be a string; change "%v" in the "%s" module for a string value`
	wrongPathMsg              = `could not find the "%s" module path`
	reqFailedOnModuleGetMsg   = `could not process requirements of the "%s" module that is based on the "%s" module when getting the "%s" module`
	reqFailedOnCommandsGetMsg = `could not process requirements of the "%s" module that is based on the "%s" module when getting the "%s" module commands`
	reqFailedOnBuildResultMsg = `could not process requirements of the "%s" module that is based on the "%s" module when getting the build results path`
	reqFailedOnCopyMsg        = `could not process requirements of the "%s" module that is based on the "%s" module when copying artifacts`

	locFailedMsg    = `could not provide modules when initializing the location`
	circularDepsMsg = `circular dependency found between modules "%s" and "%s"`
)
