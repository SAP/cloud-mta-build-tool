package dir

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// MtaLocationParameters -MTA tool file properties
type MtaLocationParameters struct {
	// SourcePath - Path to MTA project
	SourcePath string
	// TargetPath - Path to MTA tool results
	TargetPath string
	// MtaFilename - MTA yaml filename "mta.yaml" by default
	MtaFilename string
	// IsDeploymentDescriptor - indicator of deployment descriptor usage (mtad.yaml)
	Descriptor string
}

// GetSource -Get Processed Project Path
// If not provided - current directory
func (ep *MtaLocationParameters) GetSource() string {
	if ep.SourcePath == "" {
		// TODO handle error
		p, _ := os.Getwd()
		return p
	}
	return ep.SourcePath
}

// GetTarget -Get Target Path
// If not provided - path of processed project
func (ep *MtaLocationParameters) GetTarget() string {
	if ep.TargetPath == "" {
		return ep.GetSource()
	}
	return ep.TargetPath
}

// GetTargetTmpDir -Get Target Temporary Directory path
// Subdirectory in target folder named as source project folder
func (ep *MtaLocationParameters) GetTargetTmpDir() string {
	_, file := filepath.Split(ep.GetSource())
	// append to the currentPath the file name
	return filepath.Join(ep.GetTarget(), file)
}

// GetTargetModuleDir -Get path to the packed module directory
// Subdirectory in Target Temporary Directory named by module name
func (ep *MtaLocationParameters) GetTargetModuleDir(moduleName string) string {
	return filepath.Join(ep.GetTargetTmpDir(), moduleName)
}

// GetTargetModuleZipPath -Get path to the packed module data.zip
// Subdirectory in Target Temporary Directory named by module name
func (ep *MtaLocationParameters) GetTargetModuleZipPath(moduleName string) string {
	return filepath.Join(ep.GetTargetModuleDir(moduleName), "data.zip")
}

// GetSourceModuleDir -Get path to module to be packed
// Subdirectory in Source
func (ep *MtaLocationParameters) GetSourceModuleDir(modulePath string) string {
	return filepath.Join(ep.GetSource(), modulePath)
}

// GetMtaYamlFilename -Get MTA yaml File name
func (ep *MtaLocationParameters) GetMtaYamlFilename() string {
	if ep.MtaFilename == "" {
		return "mta.yaml"
	}
	return ep.MtaFilename
}

// GetMtaYamlPath -Get MTA yaml File path
func (ep *MtaLocationParameters) GetMtaYamlPath() string {
	return filepath.Join(ep.GetSource(), ep.GetMtaYamlFilename())
}

// GetMetaPath -Get path to generated META-INF directory
func (ep *MtaLocationParameters) GetMetaPath() string {
	return filepath.Join(ep.GetTargetTmpDir(), "META-INF")
}

// GetMtadPath -Get path to generated MTAD file
func (ep *MtaLocationParameters) GetMtadPath() string {
	return filepath.Join(ep.GetMetaPath(), "mtad.yaml")
}

// GetManifestPath -Get path to generated manifest file
func (ep *MtaLocationParameters) GetManifestPath() string {
	return filepath.Join(ep.GetMetaPath(), "MANIFEST.MF")
}

// GetRelativePath - -Remove the basePath from the fullPath and get only the relative
func GetRelativePath(fullPath, basePath string) string {
	return strings.TrimPrefix(fullPath, basePath)
}

// ValidateDeploymentDescriptor -Validates Deployment Descriptor
func ValidateDeploymentDescriptor(descriptor string) error {
	if descriptor != "" && descriptor != "dev" && descriptor != "dep" {
		return errors.New("Wrong descriptor value. Expected one of [dev, dep]. Default is dev")
	}
	return nil
}

// IsDeploymentDescriptor - Check if flag is related to deployment descriptor
func (ep *MtaLocationParameters) IsDeploymentDescriptor() bool {
	return ep.Descriptor == "dep"
}
