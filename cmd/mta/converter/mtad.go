package converter

import (
	"cloud-mta-build-tool/cmd/mta/models"
	"cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/constants"
	"cloud-mta-build-tool/cmd/platform"
)

// ConvertTypes file according to the deployed env
func ConvertTypes(mta models.MTA) {
	const (
		SCHEMA_VERSION = "3.1"
	)
	platformFile := dir.Load(dir.GetPath() + constants.PathSep + "platform_cfg.yaml")
	platform.Parse(platformFile)

	// Modify schema version
	*mta.SchemaVersion = SCHEMA_VERSION
	// Modify Types
	for i, element := range mta.Modules {
		element.Path = ""
		element.BuildParams = nil
		switch element.Type {
		case "html5", "sitecontent":
			mta.Modules[i].Type = "javascript.nodejs"
		case "hdb":
			mta.Modules[i].Type = "com.sap.xs.hdi"
		case "nodejs":
			mta.Modules[i].Type = "javascript.nodejs"
		case "java":
			mta.Modules[i].Type = "java.tomcat"
		}

	}
}
