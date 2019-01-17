package dir

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta/mta"
)

const (
	//Dep - deployment descriptor
	Dep = "dep"
	//Dev - development descriptor
	Dev  = "dev"
	mtad = "mtad.yaml"
)

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

// ISourceModule - source module interface
type ISourceModule interface {
	GetSourceModuleDir(modulePath string) string
}

// IMtaYaml - MTA Yaml interface
type IMtaYaml interface {
	GetMtaYamlFilename() string
	GetMtaYamlPath() string
}

// IMtaExtYaml - MTA Extension Yaml interface
type IMtaExtYaml interface {
	GetMtaExtYamlPath(platform string) string
}

// ITargetPath - target path interface
type ITargetPath interface {
	GetTarget() string
	GetTargetTmpDir() string
}

// ITargetModule - Target Module interface
type ITargetModule interface {
	GetTargetModuleDir(moduleName string) string
	GetTargetModuleZipPath(moduleName string) string
}

// IModule - module interface
type IModule interface {
	ISourceModule
	ITargetModule
}

// ITargetArtifacts - target artifacts interface
type ITargetArtifacts interface {
	GetMetaPath() string
	GetMtadPath() string
	GetManifestPath() string
}

// Loc - MTA tool file properties
type Loc struct {
	// SourcePath - Path to MTA project
	SourcePath string
	// TargetPath - Path to results
	TargetPath string
	// MtaFilename - MTA yaml filename "mta.yaml" by default
	MtaFilename string
	// Descriptor - indicator of deployment descriptor usage (mtad.yaml)
	Descriptor string
}

// GetSource gets the processed project path;
// if it is not provided, use the current directory.
func (ep *Loc) GetSource() string {
	return ep.SourcePath
}

// GetDescriptor - gets descriptor type of Location
func (ep *Loc) GetDescriptor() string {
	if ep.Descriptor == "" {
		return Dev
	}

	return ep.Descriptor
}

// GetTarget gets the target path;
// if it is not provided, use the path of the processed project.
func (ep *Loc) GetTarget() string {
	return ep.TargetPath
}

// GetTargetTmpDir gets the temporary target directory path.
// The subdirectory in the target folder is named as the source project folder.
func (ep *Loc) GetTargetTmpDir() string {
	source := ep.GetSource()
	_, file := filepath.Split(source)
	target := ep.GetTarget()
	// append to the currentPath the file name
	return filepath.Join(target, file)
}

// GetTargetModuleDir gets the path to the packed module directory.
// The subdirectory in the temporary target directory is named by the module name.
func (ep *Loc) GetTargetModuleDir(moduleName string) string {
	dir := ep.GetTargetTmpDir()

	return filepath.Join(dir, moduleName)
}

// GetTargetModuleZipPath gets the path to the packed module data.zip file.
// The subdirectory in temporary target directory is named by the module name.
func (ep *Loc) GetTargetModuleZipPath(moduleName string) string {
	return filepath.Join(ep.GetTargetModuleDir(moduleName), "data.zip")
}

// GetSourceModuleDir gets the path to the module to be packed.
// The subdirectory is in the source.
func (ep *Loc) GetSourceModuleDir(modulePath string) string {
	return filepath.Join(ep.GetSource(), filepath.Clean(modulePath))
}

// GetMtaYamlFilename - Gets the MTA .yaml file name.
func (ep *Loc) GetMtaYamlFilename() string {
	if ep.MtaFilename == "" {
		if ep.Descriptor == Dep {
			return mtad
		}
		return "mta.yaml"
	}
	return ep.MtaFilename
}

// GetMtaYamlPath gets the MTA .yaml file path.
func (ep *Loc) GetMtaYamlPath() string {
	return filepath.Join(ep.GetSource(), ep.GetMtaYamlFilename())
}

// GetMtaExtYamlPath gets the MTA extension .yaml file path.
func (ep *Loc) GetMtaExtYamlPath(platform string) string {
	return filepath.Join(ep.GetSource(), platform+"-mtaext.yaml")
}

// GetMetaPath gets the path to the generated META-INF directory.
func (ep *Loc) GetMetaPath() string {
	return filepath.Join(ep.GetTargetTmpDir(), "META-INF")
}

// GetMtadPath gets the path to the generated MTAD file.
func (ep *Loc) GetMtadPath() string {
	return filepath.Join(ep.GetMetaPath(), mtad)
}

// GetManifestPath gets the path to the generated manifest file.
func (ep *Loc) GetManifestPath() string {
	return filepath.Join(ep.GetMetaPath(), "MANIFEST.MF")
}

// ValidateDeploymentDescriptor validates the deployment descriptor.
func ValidateDeploymentDescriptor(descriptor string) error {
	if descriptor != "" && descriptor != Dev && descriptor != Dep {
		return fmt.Errorf("the %v descriptor is invalid; expected one of the following values: Dev, Dep", descriptor)
	}
	return nil
}

// IsDeploymentDescriptor checks whether the flag is related to the deployment descriptor.
func (ep *Loc) IsDeploymentDescriptor() bool {
	return ep.Descriptor == Dep
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
		// extension is not mandatory
		return &mta.EXT{}, nil
	}
	// Parse MTA extension file
	return mta.UnmarshalExt(yamlContent)
}

// Location - provides Location parameters of MTA
func Location(source, target, descriptor string, wdGetter func() (string, error)) (*Loc, error) {

	err := ValidateDeploymentDescriptor(descriptor)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize the location when validating descriptor")
	}

	var mtaFilename string
	if descriptor == Dev || descriptor == "" {
		mtaFilename = "mta.yaml"
		descriptor = Dev
	} else {
		mtaFilename = "mtad.yaml"
		descriptor = Dep
	}

	if source == "" {
		source, err = wdGetter()
		if err != nil {
			return nil, errors.Wrap(err, "failed to initialize the location when getting working directory")
		}
	}
	if target == "" {
		target = source
	}
	return &Loc{SourcePath: source, TargetPath: target, MtaFilename: mtaFilename, Descriptor: descriptor}, nil
}
