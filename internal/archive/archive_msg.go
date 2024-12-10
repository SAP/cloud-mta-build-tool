package dir

const (
	// FolderCreationFailedMsg - message raised when folder creation fails because of the existence file with identical name
	FolderCreationFailedMsg = `could not create the "%s" folder because a file exists with the same name`

	copyFailedOnGetStatusMsg         = `could not copy files matching the "%s" pattern from the "%s" folder to the "%s" folder: could not get the status of the "%s" file or folder`
	copyFailedMsg                    = `could not copy files matching the "%s" pattern from the "%s" folder to the "%s" folder: could not copy "%s" to "%s"`
	archivingFailedOnCreateFolderMsg = `could not create the "%s" folder`
	fileCreationFailedMsg            = `could not create the "%s"" file`
	copyByPatternFailedOnCreateMsg   = `could not copy files matching the patterns [%s,...] from the "%s" folder to the "%s" folder: could not create the "%s" folder`
	copyByPatternFailedOnTargetMsg   = `could not copy files matching the patterns [%s,...] from the "%s" folder to the "%s" folder: "%s" is not a folder`
	copyByPatternFailedOnMatchMsg    = `could not copy files matching the "%s" pattern from the "%s" folder to the "%s": could not get list of files matching the "%s" pattern`
	wrongPathMsg                     = `could not find the "%s" path`

	// InitLocFailedOnWorkDirMsg - message raised on getting working directory when initializing location
	InitLocFailedOnWorkDirMsg = `could not get working directory`

	InvalidMtaYamlFilenameMsg = `the "%s" is not a valid mta yaml file name;`

	// InvalidDescMsg - invalid descriptor
	InvalidDescMsg = `the "%s" descriptor is invalid; expected one of the following values: Dev, Dep`

	copyByPatternMsg    = `copying files matching the [%s,...] patterns from the "%s" folder to the "%s" folder`
	skipSymbolicLinkMsg = `copying files from the "%s" folder to the "%s" folder: skipped the "%s" entry because it's a symbolic link`

	folderCreatedMsg = `the "%s" folder has been created`

	recursiveSymLinkMsg = `the "%s" symbolic path is recursive`
	badSymLink          = `could not read the "%s" symbolic link`
)
