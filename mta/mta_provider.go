package mta

import (
	"github.com/pkg/errors"
)

// ReadMta reads the MTA file according to it's name and path
// returns a reference to the MTA object
func ReadMta(path, filename string) (*MTA, error) {
	var mta *MTA
	yamlContent, err := ReadMtaContent(path, filename)
	// Read MTA file
	if err == nil {
		mta, err = ParseToMta(yamlContent)
	}
	return mta, err
}

// ReadMtaContent reads the MTA file according to it's name and path
// returns a []byte object represents the content of the MTA file
func ReadMtaContent(path, filename string) ([]byte, error) {
	s := &Source{Path: path, Filename: filename}
	yamlContent, err := s.Readfile()
	// Read MTA file
	if err != nil {
		err = errors.New("Error reading the MTA file: " + err.Error())
	}
	return yamlContent, err
}

// ParseToMta - Parse MTA file ([]byte object) to an MTA object
func ParseToMta(content []byte) (*MTA, error) {
	mta := &MTA{}
	// Parse MTA file
	err := mta.Parse(content)
	if err != nil {
		err = errors.New("Error parsing the MTA: " + err.Error())
	}
	return mta, err
}
