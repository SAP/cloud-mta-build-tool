package artifacts

const (
	assemblingMsg = `assembling the MTA project...`

	assemblyFailedOnCopyMsg    = `could not copy MTA artifacts to assemble`
	assemblyFailedOnMetaMsg    = `could not generate the MTA metadata`
	assemblyFailedOnMtarMsg    = `could not create the MTA archive`
	assemblyFailedOnCleanupMsg = `could not clean temporary files`

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

	genMetaMsg           = `could not generate metadata`
	genMetaPopulatingMsg = `could not generate metadata when populating the manifest file`
	genMetaMTADMsg       = `could not generate metadata when generating the MTAD file`

	genMTADTypeTypeCnvMsg = `could not generate the MTAD file when converting types according to the "%s" platform`
	genMTADMarshMsg       = `could not generate the MTAD file when marshalling the MTAD object`
	genMTADWriteMsg       = `could not generate the MTAD file when writing`

	genMTARParsingMsg = `could not generate the MTA archive`
	genMTARArchMsg    = `could not generate the MTA archive when archiving`

	buildMsg                 = `building the "%s" module...`
	buildFailedMsg           = `could not build the "%s" module`
	buildFailedOnCommandsMsg = `could not get commands for the "%s" module`
	buildFailedOnDepsMsg     = `could not process dependencies for the "%s" module`
	buildResultMsg           = `the build results of the "%s" module will be packed and saved in the "%s" folder`

	packMsg                       = `packing the "%s" module...`
	packFailedOnLocMsg            = `could not pack the "%s" module when initializing the location`
	packFailedOnCommandsMsg       = `could not pack the "%s" module when getting commands`
	packFailedOnBuildArtifactMsg  = `could not pack the "%s" module while getting the build artifact`
	packFailedOnTargetArtifactMsg = `could not pack the "%s" module while getting the build artifact target path`
	packFailedOnFolderCreationMsg = `could not pack the "%s" module when creating the "%s" folder`
	packFailedOnCopyMsg           = `could not pack the "%s" module when copying the "%s" path to the "%s" path`
	// PackFailedOnArchMsg - message raised when pack fails during archiving the module
	PackFailedOnArchMsg = `could not pack the "%s" module when archiving`

	copyContentFailedOnLocMsg = `could not copy the MTA content when initializing the deployment descriptor location`
	copyContentFailedMsg      = `could not copy the MTA content`
	pathNotExistsMsg          = `the "%s" path does not exist in the MTA project location`
	copyContentCopyFailedMsg  = `could not copy the "%s" MTA content to the "%s" target directory because: %s`
	copyStartMsg              = `copying the MTA content...`
	copyDoneMsg               = `copied "%s"`
	cleanupFailedMsg          = `could not clean up`

	invalidPlatformMsg = `invalid target platform "%s"; supported platforms are: "cf", "neo", "xsa"`
	adaptationMsg      = `could not adapt the "%s" module path property`

	// UnsupportedPhaseMsg - message raised when phase of mta project build is wrong
	UnsupportedPhaseMsg     = `the "%s" phase of MTA project build is invalid; supported phases: "pre", "post"`
	execFailedMsg           = `could not build the MTA project`
	removeFailedMsg         = `could not remove the "%s" file`
	commandsMissingMsg      = `the "commands" property is missing in the "custom" builder`
	commandsNotSupportedMsg = `the "commands" property is not supported by the "%s" builder`

	validationMsg             = "validating the MTA project"
	wrongStrictnessMsg        = `the "%s" strictness value is wrong; boolean value expected`
	validationFailedOnLocMsg  = `could not validate when initializing the location`
	validationFailedOnModeMsg = `could not validate when analyzing the validation mode`

	mergeInfoMsg                 = `merging the "mta.yaml" file with the MTA extension descriptors...`
	mergeNameRequiredMsg         = `the mandatory parameter target-file-name is missing`
	mergeFailedOnFileCreationMsg = `could not create the "%s" file because it already exists`
)
