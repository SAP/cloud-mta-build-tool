package constants

import "os"

const (
	// PathSep - os path
	PathSep = string(os.PathSeparator)
	// DataZip - zip suffix
	DataZip = "/data.zip"
	// MtaYaml - mta.yaml file name
	MtaYaml = "mta.yaml"
	// MtarSuffix - mtar suffix
	MtarSuffix = ".mtar"
	// TempFolder - prefix
	TempFolder = "BUILD_MTAR_TEMP"
)
