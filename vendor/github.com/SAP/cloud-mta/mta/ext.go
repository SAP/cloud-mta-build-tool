package mta

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// UnmarshalExt returns a reference to the EXT object from a byte array.
func UnmarshalExt(content []byte) (*EXT, error) {
	m := &EXT{}
	// Unmarshal MTA file
	err := yaml.Unmarshal([]byte(content), &m)
	if err != nil {
		err = errors.Wrap(err, "failed to unmarshal the mta extension object")
	}
	return m, err
}

// Merge merges mta object with mta extension object
// extension properties complement and overwrite mta properties
func Merge(mta *MTA, mtaExt *EXT) {
	for _, module := range mta.Modules {
		extModule, err := mtaExt.GetModuleByName(module.Name)
		if err == nil {
			extendMap(&module.Properties, &extModule.Properties)
			extendMap(&module.Parameters, &extModule.Parameters)
			extendMap(&module.BuildParams, &extModule.BuildParams)
		}
	}
}

// extendMap extends map with elements of mta extension map
func extendMap(m *map[string]interface{}, ext *map[string]interface{}) {
	if *m == nil {
		*m = make(map[string]interface{})
	}
	if ext != nil {
		for key, value := range *ext {
			(*m)[key] = value
		}
	}
}
