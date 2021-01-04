package commands

import "github.com/SAP/cloud-mta-build-tool/internal/logs"

// List of deprecated builders
var deprecatedBuilders = map[string]string{"maven_deprecated": `the "maven_deprecated" builder is deprecated and will be removed after July, 2021.`}

func awareOfDeprecatedBuilder(builder string) {
	warn := deprecatedBuilders[builder]
	if warn != "" {
		logs.Logger.Warn(warn)
	}
}
