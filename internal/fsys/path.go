package dir

import (
	"os"
	"path/filepath"
	"strings"
)

// EndPoints - source and target paths of
type EndPoints struct {
	SourcePath  string
	TargetPath  string
	MtaFilename string
}

// Get Processed Project Path
// If not provided - current directory
func (ep EndPoints) GetSource() string {
	if ep.SourcePath == "" {
		p, _ := os.Getwd()
		return p
	} else {
		return ep.SourcePath
	}
}

// Get Target Path
// If not provided - path of processed project
func (ep EndPoints) GetTarget() string {
	if ep.TargetPath == "" {
		return ep.GetSource()
	} else {
		return ep.TargetPath
	}
}

// Get Target Temporary Directory path
// Subdirectory in target folder named as source project folder
func (ep EndPoints) GetTargetTmpDir() string {
	_, file := filepath.Split(ep.GetSource())
	// append to the currentPath the file name
	return filepath.Join(ep.GetTarget(), file)
}

// Get path to the packed module directory
// Subdirectory in Target Temporary Directory named by module name
func (ep EndPoints) GetTargetModuleDir(moduleName string) string {
	return filepath.Join(ep.GetTargetTmpDir(), moduleName)
}

// Get path to the packed module data.zip
// Subdirectory in Target Temporary Directory named by module name
func (ep EndPoints) GetTargetModuleZipPath(moduleName string) string {
	return filepath.Join(ep.GetTargetModuleDir(moduleName), "data.zip")
}

// Get path to module to be packed
// Subdirectory in Source
func (ep EndPoints) GetSourceModuleDir(modulePath string) string {
	return filepath.Join(ep.GetSource(), modulePath)
}

func (ep EndPoints) GetMtaYamlFilename() string {
	if ep.MtaFilename == "" {
		return "mta.yaml"
	} else {
		return ep.MtaFilename
	}
}

func (ep EndPoints) GetMtaYamlPath() string {
	return filepath.Join(ep.GetSource(), ep.GetMtaYamlFilename())
}

func (ep EndPoints) GetMetaPath() string {
	return filepath.Join(ep.GetTargetTmpDir(), "META-INF")
}

func (ep EndPoints) GetMtadPath() string {
	return filepath.Join(ep.GetMetaPath(), "mtad.yaml")
}

func (ep EndPoints) GetManifestPath() string {
	return filepath.Join(ep.GetMetaPath(), "MANIFEST.MF")
}

// GetRelativePath - remove the basePath from the fullPath and get only the relative
func GetRelativePath(fullPath, basePath string) string {
	return strings.TrimPrefix(fullPath, basePath)
}
