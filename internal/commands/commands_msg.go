package commands

const (
	missingPropMsg           = `the "commands" property is missing in the "custom" builder`
	wrongPropMsg             = `the "commands" property is defined incorrectly; the property must contain a sequence of strings`
	parseModuleCfgFailedMsg  = `could not parse the module types configuration`
	parseBuilderCfgFailedMsg = `could not parse the builder types configuration`
	wrongModuleTypeDefMsg    = `the module type definition can include either the builder or the commands; the %s module type includes both`
	undefinedBuilderMsg      = `the "%s" builder is not defined in the custom commands configuration`
	undefinedModuleMsg       = `the "%s" module is not defined in the MTA file`
	// BadCommandMsg is an error message that is returned when a command line cannot be parsed
	BadCommandMsg          = `could not parse command "%s"`
	notNativeBuilderMsg    = `the "%s" builder is not a natvie builder`
	notNativeModuleTypeMsg = `the "%s" type is not a native module type`
	emptySBomFileInputMsg  = `no sbom files in tmp dir to merge`
)
