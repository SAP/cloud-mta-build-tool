package converter

import "mbtv2/cmd/mta/models"

// ModifyMtad file according to the deployed env
func ModifyMtad(mta models.MTA) {

	const (
		SCHEMA_VERSION = "3.1"
	)

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
