package platform

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/SAP/cloud-mta/mta"
)

// Unmarshal - unmarshal platform config
func Unmarshal(data []byte) (Platforms, error) {
	platforms := Platforms{}
	err := yaml.UnmarshalStrict(data, &platforms)
	if err != nil {
		return platforms, errors.Wrap(err, UnmarshalFailedMsg)
	}
	return platforms, nil
}

// ConvertTypes - convert schema type
func ConvertTypes(iCfg mta.MTA, eCfg Platforms, targetPlatform string) {
	tpl := platformConfig(eCfg, targetPlatform)
	for i, v := range iCfg.Modules {
		moduleAcc := -1
		modulePlatformType := v.Type
		for _, em := range tpl.Modules {
			if ok, acc := satisfiesModuleConfig(v, &em); ok && acc > moduleAcc {
				modulePlatformType = em.PlatformType
				moduleAcc = acc
			}
		}
		iCfg.Modules[i].Type = modulePlatformType
	}
}

// Satisfies checks if the module m satisfies the conditions defined in the configuration mc.
//
// If it doesn't satisfy the conditions, ok will be false and accuracy will be less than 0.
//
// If it satisfies the conditions, accuracy will be higher the more conditions there are inside the configuration
// (in other words, the more specific match is considered more accurate).
func satisfiesModuleConfig(m *mta.Module, mc *Properties) (ok bool, accuracy int) {
	if m.Type != mc.NativeType {
		return false, -1
	}
	for ckey, cval := range mc.Parameters {
		if mval, ok := m.Parameters[ckey]; !ok || mval != cval {
			return false, -1
		}
	}
	for ckey, cval := range mc.Properties {
		if mval, ok := m.Properties[ckey]; !ok || mval != cval {
			return false, -1
		}
	}
	return true, len(mc.Parameters) + len(mc.Properties)
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
