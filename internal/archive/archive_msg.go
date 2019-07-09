package dir

const (
	// FolderCreationFailedMsg - message raised when folder creation fails because of the existence file with identical name
	FolderCreationFailedMsg = `could not create the "%s" folder because a file exists with the same name`

	copyFailedOnGetStatusMsg         = `could not copy the "%s" pattern from the "%s" folder to the "%s" folder when getting the status of the source entry: %s`
	copyFailedMsg                    = `could not copy the "%s" pattern from the "%s" folder to the "%s" folder when copying the "%s" entry to the "%s" entry`
	archivingFailedOnCreateFolderMsg = `could not archive when creating the "%s" folder`
	fileCreationFailedMsg            = `could not create of the "%s"" file failed`
	copyByPatternFailedOnCreateMsg   = `could not copy the patterns [%s,...] from the "%s" folder to the "%s" folder when creating the target folder`
	copyByPatternFailedOnTargetMsg   = `could not copy the patterns [%s,...] from the "%s" folder to the "%s" folder because the target is not a folder`
	copyByPatternFailedOnMatchMsg    = `could not copy the "%s" pattern from the "%s" folder to the "%s" folder when getting matching entries`

	// InitLocFailedOnDescMsg - message raised when location initialization failed on descriptor validation
	InitLocFailedOnDescMsg = "could not initialize the location when validating descriptor"

	// InitLocFailedOnWorkDirMsg - message raised on getting working directory when initializing location
	InitLocFailedOnWorkDirMsg = "could not initialize the location when getting working directory"
	invalidDescMsg            = `the "%s" descriptor is invalid; expected one of the following values: Dev, Dep`

	copyByPatternMsg    = "copying the patterns [%s,...] from the %s folder to the %s folder"
	skipSymbolicLinkMsg = `copying of the entries from the "%s" folder to the "%s" folder skipped the "%s" entry because its mode is a symbolic link`

	// ReadFailedMsg - read failed message
	ReadFailedMsg = `could not read the "%s" file`
)
