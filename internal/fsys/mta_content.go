package dir

import (
	"cloud-mta-build-tool/mta"
)

// ParseFile returns a reference to the MTA object from a given mta.yaml file.
func ParseFile(ep *Loc) (*mta.MTA, error) {
	yamlContent, err := Read(ep)
	if err != nil {
		return nil, err
	}
	// ParseFile MTA file
	return mta.Unmarshal(yamlContent)
}
