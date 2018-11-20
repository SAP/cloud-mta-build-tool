package mta

import (
	"github.com/pkg/errors"
)

// ReadMta reads the MTA file according to its name and path and
// returns a reference to the MTA object.
func ReadMta(ep *MtaLocationParameters) (*MTA, error) {
	var mta *MTA
	yamlContent, err := ReadMtaContent(ep)
	// Read MTA file
	if err == nil {
		mta, err = ParseToMta(yamlContent)
	}
	return mta, err
}

// ReadMtaContent returns a []byte array that represents the content of the MTA file
// according to the name and path.
func ReadMtaContent(ep *MtaLocationParameters) ([]byte, error) {
	yamlContent, err := ReadMtaYaml(ep)
	// Read MTA file
	if err != nil {
		err = errors.Wrap(err, "Error reading the MTA file")
	}
	return yamlContent, err
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
