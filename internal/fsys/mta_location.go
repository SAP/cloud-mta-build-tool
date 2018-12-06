package dir

import (
	"os"
	"path/filepath"

	"cloud-mta-build-tool/mta"

	"github.com/pkg/errors"
)

const (
	dep  = "dep"
	mtad = "mtad.yaml"
)

// ILoc - Location interface
type ILoc interface {
	ISource
	ITarget
	IMtaParser
}

// IMtaParser - MTA Parser interface
type IMtaParser interface {
	ParseFile() (*mta.MTA, error)
	ParseExtFile(platform string) (*mta.EXT, error)
}

// IDescriptor - descriptor interface
type IDescriptor interface {
	IsDeploymentDescriptor() bool
	GetDescriptor() string
}

// ISource - source interface
type ISource interface {
	GetSource() (string, error)
	ISourceModule
	ISourceArtifacts
}

// ISourceModule - source module interface
type ISourceModule interface {
	GetSourceModuleDir(modulePath string) (string, error)
}

// ISourceArtifacts - source artifacts interface
type ISourceArtifacts interface {
	IMtaYaml
	IMtaExtYaml
	IDescriptor
}

// IMtaYaml - MTA Yaml interface
type IMtaYaml interface {
	GetMtaYamlFilename() string
	GetMtaYamlPath() (string, error)
}

// IMtaExtYaml - MTA Extension Yaml interface
type IMtaExtYaml interface {
	GetMtaExtYamlPath(platform string) (string, error)
}

// ITarget - target interface
type ITarget interface {
	ITargetPath
	ITargetArtifacts
	ITargetModule
}

// ITargetPath - target path interface
type ITargetPath interface {
	GetTarget() (string, error)
	GetTargetTmpDir() (string, error)
}

// ITargetModule - Target Module interface
type ITargetModule interface {
	GetTargetModuleDir(moduleName string) (string, error)
	GetTargetModuleZipPath(moduleName string) (string, error)
}

// IModule - module interface
type IModule interface {
	ISourceModule
	ITargetModule
}

// ITargetArtifacts - target artifacts interface
type ITargetArtifacts interface {
	GetMetaPath() (string, error)
	GetMtadPath() (string, error)
	GetManifestPath() (string, error)
}

// Loc - MTA tool file properties
type Loc struct {
	// SourcePath - Path to MTA project
	SourcePath string
	// TargetPath - Path to results
	TargetPath string
	// MtaFilename - MTA yaml filename "mta.yaml" by default
	MtaFilename string
	// IsDeploymentDescriptor - indicator of deployment descriptor usage (mtad.yaml)
	Descriptor string
}

// osGetWd - get working dir
var osGetWd = func() (string, error) {
	return os.Getwd()
}

// getWorkingDirectory assignment
var getWorkingDirectory = osGetWd

// GetSource gets the processed project path;
// if it is not provided, use the current directory.
func (ep *Loc) GetSource() (string, error) {
	if ep.SourcePath == "" {
		wd, err := getWorkingDirectory()
		if err != nil {
			return "", errors.Wrap(err, "GetSource failed")
		}
		return wd, nil
	}
	return ep.SourcePath, nil
}

// GetDescriptor - gets descriptor type of location
func (ep *Loc) GetDescriptor() string {
	if ep.Descriptor == "" {
		return "dev"
	}

	return ep.Descriptor
}

// GetTarget gets the target path;
// if it is not provided, use the path of the processed project.
func (ep *Loc) GetTarget() (string, error) {
	if ep.TargetPath == "" {
		source, err := ep.GetSource()
		if err != nil {
			return "", errors.Wrap(err, "GetTarget failed")
		}
		return source, nil
	}
	return ep.TargetPath, nil
}

// GetTargetTmpDir gets the temporary target directory path.
// The subdirectory in the target folder is named as the source project folder.
func (ep *Loc) GetTargetTmpDir() (string, error) {
	source, err := ep.GetSource()
	if err != nil {
		return "", errors.Wrap(err, "GetTargetTmpDir failed")
	}
	_, file := filepath.Split(source)
	target, err := ep.GetTarget()
	if err != nil {
		return "", errors.Wrap(err, "GetTargetTmpDir failed")
	}
	// append to the currentPath the file name
	return filepath.Join(target, file), nil
}

// GetTargetModuleDir gets the path to the packed module directory.
// The subdirectory in the temporary target directory is named by the module name.
func (ep *Loc) GetTargetModuleDir(moduleName string) (string, error) {
	dir, err := ep.GetTargetTmpDir()
	if err != nil {
		return "", errors.Wrap(err, "GetTargetModuleDir failed")
	}

	return filepath.Join(dir, moduleName), nil
}

// GetTargetModuleZipPath gets the path to the packed module data.zip file.
// The subdirectory in temporary target directory is named by the module name.
func (ep *Loc) GetTargetModuleZipPath(moduleName string) (string, error) {
	dir, err := ep.GetTargetModuleDir(moduleName)
	if err != nil {
		return "", errors.Wrap(err, "GetTargetModuleZipPath failed")
	}
	return filepath.Join(dir, "data.zip"), nil
}

// GetSourceModuleDir gets the path to the module to be packed.
// The subdirectory is in the source.
func (ep *Loc) GetSourceModuleDir(modulePath string) (string, error) {
	source, err := ep.GetSource()
	if err != nil {
		return "", errors.Wrap(err, "GetSourceModuleDir failed")
	}
	return filepath.Join(source, filepath.Clean(modulePath)), nil
}

// GetMtaYamlFilename - Gets the MTA .yaml file name.
func (ep *Loc) GetMtaYamlFilename() string {
	if ep.MtaFilename == "" {
		if ep.Descriptor == dep {
			return mtad
		}
		return "mta.yaml"
	}
	return ep.MtaFilename
}

// GetMtaYamlPath gets the MTA .yaml file path.
func (ep *Loc) GetMtaYamlPath() (string, error) {
	source, err := ep.GetSource()
	if err != nil {
		return "", errors.Wrap(err, "GetMtaYamlPath failed")
	}
	return filepath.Join(source, ep.GetMtaYamlFilename()), nil
}

// GetMtaExtYamlPath gets the MTA extension .yaml file path.
func (ep *Loc) GetMtaExtYamlPath(platform string) (string, error) {
	source, err := ep.GetSource()
	if err != nil {
		return "", errors.Wrap(err, "GetMtaExtYamlPath failed")
	}
	return filepath.Join(source, platform+"-mtaext.yaml"), nil
}

// GetMetaPath gets the path to the generated META-INF directory.
func (ep *Loc) GetMetaPath() (string, error) {
	dir, err := ep.GetTargetTmpDir()
	if err != nil {
		return "", errors.Wrap(err, "GetMetaPath failed")
	}
	return filepath.Join(dir, "META-INF"), nil
}

// GetMtadPath gets the path to the generated MTAD file.
func (ep *Loc) GetMtadPath() (string, error) {
	dir, err := ep.GetMetaPath()
	if err != nil {
		return "", errors.Wrap(err, "GetMtadPath failed")
	}
	return filepath.Join(dir, mtad), nil
}

// GetManifestPath gets the path to the generated manifest file.
func (ep *Loc) GetManifestPath() (string, error) {
	dir, err := ep.GetMetaPath()
	if err != nil {
		return "", errors.Wrap(err, "GetManifestPath failed")
	}
	return filepath.Join(dir, "MANIFEST.MF"), nil
}

// ValidateDeploymentDescriptor validates the deployment descriptor.
func ValidateDeploymentDescriptor(descriptor string) error {
	if descriptor != "" && descriptor != "dev" && descriptor != dep {
		return errors.New("Wrong descriptor value. Expected one of [dev, dep]. Default is dev")
	}
	return nil
}

// IsDeploymentDescriptor checks whether the flag is related to the deployment descriptor.
func (ep *Loc) IsDeploymentDescriptor() bool {
	return ep.Descriptor == dep
}

// ParseFile returns a reference to the MTA object from a given mta.yaml file.
func (ep *Loc) ParseFile() (*mta.MTA, error) {
	yamlContent, err := Read(ep)
	if err != nil {
		return nil, err
	}
	// Parse MTA file
	return mta.Unmarshal(yamlContent)
}

// ParseExtFile returns a reference to the MTA object from a given mta.yaml file.
func (ep *Loc) ParseExtFile(platform string) (*mta.EXT, error) {
	yamlContent, err := ReadExt(ep, platform)
	if err != nil {
		return nil, err
	}
	// Parse MTA extension file
	return mta.UnmarshalExt(yamlContent)
}
