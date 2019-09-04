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

	// InitLocFailedOnWorkDirMsg - message raised on getting working directory when initializing location
	InitLocFailedOnWorkDirMsg = `could not get working directory`

	// InvalidDescMsg - invalid descriptor
	InvalidDescMsg = `the "%s" descriptor is invalid; expected one of the following values: Dev, Dep`

	copyByPatternMsg    = `copying files matching the [%s,...] patterns from the "%s" folder to the "%s" folder`
	skipSymbolicLinkMsg = `copying files from the "%s" folder to the "%s" folder: skipped the "%s" entry because it's a symbolic link`

	// ReadFailedMsg - read failed message
	ReadFailedMsg = `could not read the "%s" file`

	folderCreatedMsg = `the "%s" folder has been created`

	parseExtFileFailed = `the "%s" file is not a valid MTA extension descriptor`
	// ParseMtaYamlFileFailedMsg - parse of mta yaml file failed
	ParseMtaYamlFileFailedMsg = `the "%s" file is not a valid MTA descriptor`
	extensionIDSameAsMtaIDMsg = `the "%s" extension descriptor file has the same ID ("%s") as the "%s" file`
	duplicateExtensionIDMsg   = `more than 1 extension descriptor file ("%s", "%s", ...) has the same ID ("%s")`
	duplicateExtendsMsg       = `more than 1 extension descriptor file ("%s", "%s", ...) extends the same ID ("%s")`
	extendsMsg                = `the "%s" file extends "%s"`
	unknownExtendsMsg         = `some MTA extension descriptors extend unknown IDs: %s`
)
