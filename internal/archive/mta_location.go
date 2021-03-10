package dir

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta/mta"
)

const (
	//Dep - deployment descriptor
	Dep = "dep"
	//Dev - development descriptor
	Dev = "dev"
	// TempFolderSuffix - temporary folder suffix
	TempFolderSuffix = "_mta_build_tmp"
	// Mtad - deployment descriptor file name
	Mtad = "mtad.yaml"
	// MtarFolder - default archives folder
	MtarFolder = "mta_archives"
)

// IMtaParser - MTA Parser interface
type IMtaParser interface {
	ParseFile() (*mta.MTA, error)
}

// IDescriptor - descriptor interface
type IDescriptor interface {
	IsDeploymentDescriptor() bool
	GetDescriptor() string
}

// ISourceModule - source module interface
type ISourceModule interface {
	GetSourceModuleDir(modulePath string) string
	GetSourceModuleArtifactRelPath(modulePath, artifactPath string) (string, error)
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
	GetTargetTmpDir() string
	GetTargetTmpRoot() string
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
	GetMtarDir(targetProvided bool) string
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
	// ExtensionFileNames - list of MTA extension descriptors (could be empty)
	ExtensionFileNames []string
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

// GetMtarDir - gets archive folder
// if the target folder provided archive will be saved in the target folder
// otherwise archives folder - "mta_archives" subfolder in the project folder
func (ep *Loc) GetMtarDir(targetProvided bool) string {
	if !targetProvided {
		return filepath.Join(ep.GetTarget(), MtarFolder)
	}

	return ep.GetTarget()
}

// GetTarget gets the target path;
// if it is not provided, use the path of the processed project.
func (ep *Loc) GetTarget() string {
	return ep.TargetPath
}

// GetTargetTmpDir gets the temporary target directory path.
// The subdirectory in the target folder is named as the source project folder suffixed with "_mta_build_tmp".
// Subdirectory name is prefixed with "." as a hidden folder
func (ep *Loc) GetTargetTmpDir() string {
	source := ep.GetSource()
	_, file := filepath.Split(source)
	file = "." + file + TempFolderSuffix
	target := ep.GetTarget()
	// append to the currentPath the file name
	return filepath.Join(target, file)
}

// GetTargetTmpRoot gets the build results directory root path.
func (ep *Loc) GetTargetTmpRoot() string {
	return ep.GetTargetTmpDir()
}

// GetTargetModuleDir gets the path to the packed module directory.
// The subdirectory in the temporary target directory is named by the module name.
func (ep *Loc) GetTargetModuleDir(moduleName string) string {
	dir := ep.GetTargetTmpDir()

	return filepath.Join(dir, moduleName)
}

// GetSourceModuleDir gets the path to the module to be packed.
// The subdirectory is in the source.
func (ep *Loc) GetSourceModuleDir(modulePath string) string {
	return filepath.Join(ep.GetSource(), filepath.Clean(modulePath))
}

// GetSourceModuleArtifactRelPath gets the relative path to the module's artifact
func (ep *Loc) GetSourceModuleArtifactRelPath(moduleRelPath, artifactAbsPath string) (string, error) {
	info, err := os.Stat(artifactAbsPath)
	if err != nil {
		return "", err
	}
	isFolder := info.IsDir()
	modulePath := ep.GetSourceModuleDir(moduleRelPath)
	if isFolder {
		return filepath.Rel(modulePath, artifactAbsPath)
	} else if artifactAbsPath == modulePath {
		return "", nil
	}
	return filepath.Rel(modulePath, filepath.Dir(artifactAbsPath))
}

// GetMtaYamlFilename - Gets the MTA .yaml file name.
func (ep *Loc) GetMtaYamlFilename() string {
	if ep.MtaFilename == "" {
		if ep.Descriptor == Dep {
			return Mtad
		}
		return "mta.yaml"
	}
	return ep.MtaFilename
}

// GetMtaYamlPath gets the MTA .yaml file path.
func (ep *Loc) GetMtaYamlPath() string {
	return filepath.Join(ep.GetSource(), ep.GetMtaYamlFilename())
}

// GetMtaExtYamlPath gets the full MTA extension file path by file name or path.
// If the file name is an absolute path it's returned as is. Otherwise the returned path is relative
// to the source folder.
func (ep *Loc) GetMtaExtYamlPath(extFileName string) string {
	if filepath.IsAbs(extFileName) {
		return extFileName
	}

	return filepath.Join(ep.GetSource(), extFileName)
}

// GetMetaPath gets the path to the generated META-INF directory.
func (ep *Loc) GetMetaPath() string {
	return filepath.Join(ep.GetTargetTmpDir(), "META-INF")
}

// GetMtadPath gets the path to the generated MTAD file.
func (ep *Loc) GetMtadPath() string {
	return filepath.Join(ep.GetMetaPath(), Mtad)
}

// GetManifestPath gets the path to the generated manifest file.
func (ep *Loc) GetManifestPath() string {
	return filepath.Join(ep.GetMetaPath(), "MANIFEST.MF")
}

// ValidateDeploymentDescriptor validates the deployment descriptor.
func ValidateDeploymentDescriptor(descriptor string) error {
	if descriptor != "" && descriptor != Dev && descriptor != Dep {
		return fmt.Errorf(InvalidDescMsg, descriptor)
	}
	return nil
}

// IsDeploymentDescriptor checks whether the flag is related to the deployment descriptor.
func (ep *Loc) IsDeploymentDescriptor() bool {
	return ep.Descriptor == Dep
}

// ParseFile returns a reference to the MTA object resulting from the given mta.yaml file merged with the extension descriptors.
func (ep *Loc) ParseFile() (*mta.MTA, error) {
	mtaFile, _, err := mta.GetMtaFromFile(ep.GetMtaYamlPath(), ep.GetExtensionFilePaths(), true)
	return mtaFile, err
}

// GetExtensionFilePaths returns the MTA extension descriptor full paths
func (ep *Loc) GetExtensionFilePaths() []string {
	paths := make([]string, len(ep.ExtensionFileNames))
	for i, fileName := range ep.ExtensionFileNames {
		paths[i] = ep.GetMtaExtYamlPath(fileName)
	}
	return paths
}

// Location - provides Location parameters of MTA
func Location(source, target, descriptor string, extensions []string, wdGetter func() (string, error)) (*Loc, error) {

	err := ValidateDeploymentDescriptor(descriptor)
	if err != nil {
		return nil, err
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
			return nil, errors.Wrap(err, InitLocFailedOnWorkDirMsg)
		}
	}

	if target == "" {
		target = source
	}
	return &Loc{
		SourcePath:         filepath.Join(source),
		TargetPath:         filepath.Join(target),
		MtaFilename:        mtaFilename,
		Descriptor:         descriptor,
		ExtensionFileNames: extensions,
	}, nil
}
