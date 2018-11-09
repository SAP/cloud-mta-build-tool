package platform

import (
	"cloud-mta-build-tool/mta"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Parse - parse platform config
func Parse(data []byte) (Platforms, error) {
	platforms := Platforms{}
	err := yaml.Unmarshal(data, &platforms)
	if err != nil {
		return platforms, errors.Wrap(err, "Yaml file is not valid")
	}
	return platforms, nil
}

// ConvertTypes - convert schema type
func ConvertTypes(iCfg mta.MTA, eCfg Platforms, targetPlatform string) {
	// todo get from config
	const (
		SchemaVersion = "3.1"
	)
	tpl := platformConfig(eCfg, targetPlatform)
	for i, v := range iCfg.Modules {
		*iCfg.SchemaVersion = SchemaVersion
		for _, em := range tpl.Modules {
			if v.Type == em.NativeType {
				iCfg.Modules[i].Type = em.PlatformType
			}
		}
	}
}

func platformConfig(eCfg Platforms, targetPlatform string) Modules {
	var tpl Modules
	for _, tp := range eCfg.Platforms {
		if tp.Name == targetPlatform {
			tpl.Name = tp.Name
			tpl.Modules = tp.Modules
		}
	}
	return tpl
}
