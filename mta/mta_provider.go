package mta

import (
	"cloud-mta-build-tool/internal/fsys"
	"github.com/pkg/errors"
)

// ReadMta reads the MTA file according to it's name and path
// returns a reference to the MTA object
func ReadMta(ep dir.EndPoints) (*MTA, error) {
	var mta *MTA
	yamlContent, err := ReadMtaContent(ep)
	// Read MTA file
	if err == nil {
		mta, err = ParseToMta(yamlContent)
	}
	return mta, err
}

// ReadMtaContent reads the MTA file according to it's name and path
// returns a []byte object represents the content of the MTA file
func ReadMtaContent(ep dir.EndPoints) ([]byte, error) {
	yamlContent, err := ReadMtaYaml(ep)
	// Read MTA file
	if err != nil {
		err = errors.Wrap(err, "Error reading the MTA file")
	}
	return yamlContent, err
}

// ParseToMta - Parse MTA file ([]byte object) to an MTA object
func ParseToMta(content []byte) (*MTA, error) {
	mta := &MTA{}
	// Parse MTA file
	err := mta.Parse(content)
	if err != nil {
		err = errors.Wrap(err, "Error parsing the MTA")
	}
	return mta, err
}
