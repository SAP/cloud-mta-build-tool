package dir

import (
	"os"
	"path/filepath"
	"strings"
)

// EndPoints - MTA tool file properties
type EndPoints struct {
	// SourcePath - Path to MTA project
	SourcePath string
	// TargetPath - Path to MTA tool results
	TargetPath string
	// MtaFilename - MTA yaml filename "mta.yaml" by default
	MtaFilename string
}

// GetSource Get Processed Project Path
// If not provided - current directory
func (ep EndPoints) GetSource() string {
	if ep.SourcePath == "" {
		p, _ := os.Getwd()
		return p
	} else {
		return ep.SourcePath
	}
}

// GetTarget Get Target Path
// If not provided - path of processed project
func (ep EndPoints) GetTarget() string {
	if ep.TargetPath == "" {
		return ep.GetSource()
	} else {
		return ep.TargetPath
	}
}

// GetTargetTmpDir Get Target Temporary Directory path
// Subdirectory in target folder named as source project folder
func (ep EndPoints) GetTargetTmpDir() string {
	_, file := filepath.Split(ep.GetSource())
	// append to the currentPath the file name
	return filepath.Join(ep.GetTarget(), file)
}

// GetTargetModuleDir Get path to the packed module directory
// Subdirectory in Target Temporary Directory named by module name
func (ep EndPoints) GetTargetModuleDir(moduleName string) string {
	return filepath.Join(ep.GetTargetTmpDir(), moduleName)
}

// GetTargetModuleZipPath Get path to the packed module data.zip
// Subdirectory in Target Temporary Directory named by module name
func (ep EndPoints) GetTargetModuleZipPath(moduleName string) string {
	return filepath.Join(ep.GetTargetModuleDir(moduleName), "data.zip")
}

// GetSourceModuleDir Get path to module to be packed
// Subdirectory in Source
func (ep EndPoints) GetSourceModuleDir(modulePath string) string {
	return filepath.Join(ep.GetSource(), modulePath)
}

// GetMtaYamlFilename Get MTA yaml File name
func (ep EndPoints) GetMtaYamlFilename() string {
	if ep.MtaFilename == "" {
		return "mta.yaml"
	} else {
		return ep.MtaFilename
	}
}

// GetMtaYamlPath Get MTA yaml File path
func (ep EndPoints) GetMtaYamlPath() string {
	return filepath.Join(ep.GetSource(), ep.GetMtaYamlFilename())
}

// GetMetaPath - Get path to generated META-INF directory
func (ep EndPoints) GetMetaPath() string {
	return filepath.Join(ep.GetTargetTmpDir(), "META-INF")
}

// GetMtadPath Get path to generated MTAD file
func (ep EndPoints) GetMtadPath() string {
	return filepath.Join(ep.GetMetaPath(), "mtad.yaml")
}

// GetManifestPath Get path to generated manifest file
func (ep EndPoints) GetManifestPath() string {
	return filepath.Join(ep.GetMetaPath(), "MANIFEST.MF")
}

// GetRelativePath - remove the basePath from the fullPath and get only the relative
func GetRelativePath(fullPath, basePath string) string {
	return strings.TrimPrefix(fullPath, basePath)
}
