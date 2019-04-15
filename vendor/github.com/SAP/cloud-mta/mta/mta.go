// Package mta provides a convenient way of exploring the structure of `mta.yaml` file objects
// such as retrieving a list of resources required by a specific module.
package mta

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// GetModules returns a list of MTA modules.
func (mta *MTA) GetModules() []*Module {
	return mta.Modules
}

// GetResources returns list of MTA resources.
func (mta *MTA) GetResources() []*Resource {
	return mta.Resources
}

// GetModuleByName returns a specific module by name.
func (mta *MTA) GetModuleByName(name string) (*Module, error) {
	for _, m := range mta.Modules {
		if m.Name == name {
			return m, nil
		}
	}
	return nil, fmt.Errorf("the %s module is not defined ", name)
}

// GetModuleByName returns a specific module by name from extension object
func (ext *EXT) GetModuleByName(name string) (*ModuleExt, error) {
	for _, m := range ext.Modules {
		if m.Name == name {
			return m, nil
		}
	}
	return nil, fmt.Errorf("the %s module is not defined ", name)
}

// GetResourceByName returns a specific resource by name.
func (mta *MTA) GetResourceByName(name string) (*Resource, error) {
	for _, r := range mta.Resources {
		if r.Name == name {
			return r, nil
		}
	}
	return nil, fmt.Errorf("the %s resource is not defined ", name)
}

// Unmarshal returns a reference to the MTA object from a byte array.
func Unmarshal(content []byte) (*MTA, error) {
	m := &MTA{}
	// Unmarshal MTA file
	err := yaml.Unmarshal([]byte(content), &m)
	return m, err
}
