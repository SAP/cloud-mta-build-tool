package converter

import (
	"cloud-mta-build-tool/cmd/mta/models"
	"cloud-mta-build-tool/cmd/platform"
)

func ConvertTypes(iCfg models.MTA, eCfg platform.Platforms, targetPlatform string) {
	//todo get from config
	const (
		SchemaVersion = "3.1"
	)
	tpl := PlatformConfig(eCfg, targetPlatform)
	for i, v := range iCfg.Modules {
		*iCfg.SchemaVersion = SchemaVersion
		for _, em := range tpl.Modules {
			if v.Type == em.NativeType {
				iCfg.Modules[i].Type = em.PlatformType
			}
		}
	}
}

func PlatformConfig(eCfg platform.Platforms, targetPlatform string) platform.Modules {
	var tpl platform.Modules
	for _, tp := range eCfg.Platforms {
		if tp.Name == targetPlatform {
			tpl.Name = tp.Name
			tpl.Modules = tp.Modules
		}
	}
	return tpl
}
