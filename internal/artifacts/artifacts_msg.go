package artifacts

const (
	assemblingMsg = `assembling the MTA project...`

	assemblyFailedOnCopyMsg    = `assembly of the MTA project failed when copying the MTA content`
	assemblyFailedOnMetaMsg    = `assembly of the MTA project failed when generating the meta information`
	assemblyFailedOnMtarMsg    = `assembly of the MTA project failed when generating the MTA archive`
	assemblyFailedOnCleanupMsg = `assembly of the MTA project failed when executing cleanup`

	cleanupMsg               = `cleaning temporary files...`
	cleanupFailedOnLocMsg    = `cleanup failed when initializing the location`
	cleanupFailedOnFolderMsg = `cleanup failed when removing the "%s" folder`

	wrongArtifactPathMsg          = `could not generate the manifest file when getting the artifact path of the "%s" module`
	unknownModuleContentTypeMsg   = `could not generate the manifest file when getting the "%s" module content type`
	unknownResourceContentTypeMsg = `could not generate the manifest file when getting the "%s" resource content type`
	requiredEntriesProblemMsg     = `could not generate the manifest file when building the required entries of the "%s" module`
	contentTypeDefMsg             = `the "%s" path does not exist; the content type was not defined`
	cliVersionMsg                 = `could not generate the manifest file when getting the CLI version`
	initMsg                       = `could not generate the manifest file when initializing it`
	populationMsg                 = `could not generate the manifest file when populating the content`
	contentTypeCfgMsg             = `could not generate the manifest file when getting the content types from the configuration`

	genMetaParsingMsg    = `could not generate metadata when parsing the MTA file`
	genMetaPopulatingMsg = `could not generate metadata when populating the manifest file`
	genMetaMTADMsg       = `could not generate metadata when generating the MTAD file`

	genMTADParsingMsg     = `could not generate the MTAD file when parsing the "%s" file`
	genMTADTypeTypeCnvMsg = `could not generate the MTAD file when converting types according to the "%s" platform`
	genMTADMarshMsg       = `could not generate the MTAD file when marshalling the MTAD object`
	genMTADWriteMsg       = `could not generate the MTAD file when writing`

	genMTARParsingMsg = `could not generate the MTA archive when parsing the mta file`
	genMTARArchMsg    = `could not generate the MTA archive when archiving`

	buildMsg                 = `building the "%s" module...`
	buildFailedOnLocMsg      = `could not build the "%s" module when initializing the location`
	buildFailedOnCommandsMsg = `could not build the "%s" module when getting commands`
	buildFailedOnDepsMsg     = `could not build the "%s" module when processing dependencies`
	buildFailedOnExecCmdMsg  = `could not build the "%s" module when executing commands`
	buildResultMsg           = `the build results of the "%s" module will be packed and saved in the "%s" folder`

	packMsg                       = `packing the "%s" module...`
	packFailedOnLocMsg            = `could not pack the "%s" module when initializing the location`
	packFailedOnCommandsMsg       = `could not pack the "%s" module when getting commands`
	packFailedOnBuildArtifactMsg  = `could not pack the "%s" module while getting the build artifact`
	packFailedOnTargetArtifactMsg = `could not pack the "%s" module while getting the build artifact target path`
	packFailedOnFolderCreationMsg = `could not pack of the "%s" module when creating the "%s" folder`
	packFailedOnCopyMsg           = `could not pack of the "%s" module when copying the "%s" path to the "%s" path`
	// PackFailedOnArchMsg - message raised when pack fails during archiving the module
	PackFailedOnArchMsg           = `could not pack of the "%s" module when archiving`

	copyContentFailedOnLocMsg   = `could not copy the MTA content when initializing the deployment descriptor location`
	copyContentFailedOnParseMsg = `could not copy the MTA content when parsing the %s file`
	pathNotExistsMsg            = `the "%s" path does not exist in the MTA project location`
	copyContentFailedMsg        = `could not copy the "%s" MTA content to the "%s" target directory because: %s`
	copyStartMsg                = `copying the MTA content...`
	copyDoneMsg                 = `copied "%s"`
	cleanupFailedMsg            = `; cleanup failed`

	invalidPlatformMsg = `the invalid target platform "%s"; supported platforms are: "cf", "neo", "xsa"`
	adaptationMsg      = `could not adapt the "%s" module path property`

	execFailedMsg           = `could not execute the "%s" file`
	removeFailedMsg         = `could not remove the "%s" file`
	execAndRemoveFailedMsg  = `could not execute the "%s" file; could not remove the "%s" file`
	// UnsupportedPhaseMsg - message raised when phase of mta project build is wrong
	UnsupportedPhaseMsg     = `the "%s" phase of mta project build is invalid; supported phases: "pre", "post"`
	commandsMissingMsg      = `the "commands" property is missing in the "custom" builder`
	commandsNotSupportedMsg = `the "commands" property is not supported by the "%s" builder`

	validationMsg             = "validating the MTA project"
	wrongStrictnessMsg        = `the "%s" strictness value is wrong; boolean value expected`
	validationFailedOnLocMsg  = `could not validate when initializing the location`
	validationFailedOnModeMsg = `could not validate when analyzing the validation mode`
)
