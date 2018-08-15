package converter

import (
	"strings"

	"cloud-mta-build-tool/cmd/mta/models"
	"cloud-mta-build-tool/cmd/platform"
)

// ConvertTypes file according to the deployed env
func ConvertTypes(mta models.MTA, platforms platform.Platforms, platform string) {
	//todo get from config
	const (
		SCHEMA_VERSION = "3.1"
	)
	for _, module := range platforms.Platforms {

		if module.Name == platform {
			// Modify schema version
			*mta.SchemaVersion = SCHEMA_VERSION
			// Modify Types
			for i, value := range module.Models {
				//Check for types
				if len(mta.Modules) > i {
					if strings.Compare(value.NativeType, mta.Modules[i].Type)  == 0 {
						//Modify the module type according the platform config
						mta.Modules[i].Type = value.PlatformType
					}
				}
			}
		}
	}
}
