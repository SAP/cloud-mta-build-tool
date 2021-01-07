package commands

import "github.com/SAP/cloud-mta-build-tool/internal/logs"

// List of deprecated builders
var deprecatedBuilders = map[string]string{"maven_deprecated": `the "maven_deprecated" builder is deprecated and will be removed on July 2021`}

func checkDeprecatedBuilder(builder string) {
	warn := deprecatedBuilders[builder]
	if warn != "" {
		logs.Logger.Warn(warn)
	}
}
