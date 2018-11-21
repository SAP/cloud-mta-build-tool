package mta

import (
	"github.com/pkg/errors"
)

// ReadFile reads the MTA file according to its name and path and
// returns a reference to the MTA object.
func ReadFile(ep *Loc) (*MTA, error) {
	var mta *MTA
	yamlContent, err := Read(ep)
	// ReadFile MTA file
	if err == nil {
		mta, err = ParseToMta(yamlContent)
	}
	return mta, err
}

// ParseToMta returns a byte array of an MTA object.
func ParseToMta(content []byte) (*MTA, error) {
	mta := &MTA{}
	// Parse MTA file
	err := mta.Parse(content)
	if err != nil {
		err = errors.Wrap(err, "Error parsing the MTA")
	}
	return mta, err
}
