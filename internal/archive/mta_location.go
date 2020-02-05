package dir

import (
	"fmt"
	"path/filepath"
	"strings"

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
	GetSourceModuleArtifactRelPath(modulePath, artifactPath string, artifactFolder bool) string
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
func (ep *Loc) GetSourceModuleArtifactRelPath(moduleRelPath, artifactAbsPath string, artifactFolder bool) string {
	modulePath := ep.GetSourceModuleDir(moduleRelPath)
	if artifactFolder {
		return strings.Replace(artifactAbsPath, modulePath, "", 1)
	} else if artifactAbsPath == modulePath {
		return ""
	}
	return strings.Replace(filepath.Dir(artifactAbsPath), modulePath, "", 1)
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

// ParseMtaFile returns a reference to the MTA object from a given mta.yaml file.
func (ep *Loc) ParseMtaFile() (*mta.MTA, error) {
	yamlContent, err := Read(ep)
	if err != nil {
		return nil, err
	}
	// Parse MTA file
	mtaFile, err := mta.Unmarshal(yamlContent)
	if err != nil {
		return mtaFile, errors.Wrapf(err, ParseMtaYamlFileFailedMsg, ep.GetMtaYamlFilename())
	}
	return mtaFile, nil
}

// ParseFile returns a reference to the MTA object resulting from the given mta.yaml file merged with the extension descriptors.
func (ep *Loc) ParseFile() (*mta.MTA, error) {
	mtaFile, err := ep.ParseMtaFile()
	if err != nil {
		return mtaFile, err
	}
	extensions, err := ep.getSortedExtensions(mtaFile.ID)
	if err != nil {
		return mtaFile, err
	}
	for _, extFile := range extensions {
		// Check there is no version mismatch - the extension must have the same major.minor as the MTA
		err = checkSchemaVersionMatches(mtaFile, extFile)
		if err != nil {
			return mtaFile, err
		}

		err = mta.Merge(mtaFile, extFile)
		if err != nil {
			return mtaFile, err
		}
	}
	return mtaFile, nil
}

type extensionDetails struct {
	fileName string
	ext      *mta.EXT
}

func (ep *Loc) getSortedExtensions(mtaID string) ([]*mta.EXT, error) {
	extensionFileNames := ep.ExtensionFileNames

	// Parse all extension files and put them in a slice of extension details (the extension with the file name)
	extensions, err := parseExtensionsWithDetails(extensionFileNames, ep)
	if err != nil {
		return nil, err
	}

	// Make sure each extension has its own ID
	err = checkExtensionIDsUniqueness(extensionFileNames, extensions, mtaID, ep)
	if err != nil {
		return nil, err
	}

	// Make sure each extension extends a different ID and put them in a map of extends -> extension details
	extendsMap := make(map[string]extensionDetails, len(extensionFileNames))
	for _, details := range extensions {
		if value, ok := extendsMap[details.ext.Extends]; ok {
			return nil, errors.Errorf(duplicateExtendsMsg,
				ep.GetMtaExtYamlPath(value.fileName), ep.GetMtaExtYamlPath(details.fileName), details.ext.Extends)
		}
		extendsMap[details.ext.Extends] = details
	}

	// Verify chain of extensions and put the extensions in a slice by extends order
	return sortAndVerifyExtendsChain(extensionFileNames, mtaID, extendsMap, ep)
}

func parseExtensionsWithDetails(extensionFileNames []string, ep *Loc) ([]extensionDetails, error) {
	extensions := make([]extensionDetails, len(extensionFileNames))
	for i, extFileName := range extensionFileNames {
		extFile, err := ep.ParseExtFile(extFileName)
		if err != nil {
			return nil, err
		}
		extensions[i] = extensionDetails{extFileName, extFile}
	}
	return extensions, nil
}

func checkExtensionIDsUniqueness(extensionFileNames []string, extensions []extensionDetails, mtaID string, ep *Loc) error {
	extensionIDMap := make(map[string]extensionDetails, len(extensionFileNames))
	for _, details := range extensions {
		if details.ext.ID == mtaID {
			return errors.Errorf(extensionIDSameAsMtaIDMsg,
				ep.GetMtaExtYamlPath(details.fileName), mtaID, ep.GetMtaYamlFilename())
		}
		if value, ok := extensionIDMap[details.ext.ID]; ok {
			return errors.Errorf(duplicateExtensionIDMsg,
				ep.GetMtaExtYamlPath(value.fileName), ep.GetMtaExtYamlPath(details.fileName), details.ext.ID)
		}
		extensionIDMap[details.ext.ID] = details
	}
	return nil
}

func sortAndVerifyExtendsChain(extensionFileNames []string, mtaID string, extendsMap map[string]extensionDetails, ep IMtaExtYaml) ([]*mta.EXT, error) {
	extFiles := make([]*mta.EXT, 0, len(extensionFileNames))
	currExtends := mtaID
	value, ok := extendsMap[currExtends]
	for ok {
		extFiles = append(extFiles, value.ext)
		delete(extendsMap, currExtends)
		currExtends = value.ext.ID
		value, ok = extendsMap[currExtends]
	}
	// Check if there are extensions which extend unknown files
	if len(extendsMap) > 0 {
		// Build an error that looks like this:
		// `some MTA extension descriptors extend unknown IDs: file "myext.mtaext" extends "ext1"; file "aaa.mtaext" extends "ext2"`
		fileParts := make([]string, 0, len(extendsMap))
		for extends, details := range extendsMap {
			fileParts = append(fileParts, fmt.Sprintf(extendsMsg, ep.GetMtaExtYamlPath(details.fileName), extends))
		}
		return nil, errors.Errorf(unknownExtendsMsg, strings.Join(fileParts, `; `))
	}
	return extFiles, nil
}

func checkSchemaVersionMatches(mta *mta.MTA, ext *mta.EXT) error {
	mtaVersion := ""
	if mta.SchemaVersion != nil {
		mtaVersion = *mta.SchemaVersion
	}
	extVersion := ""
	if ext.SchemaVersion != nil {
		extVersion = *ext.SchemaVersion
	}

	if strings.SplitN(mtaVersion, ".", 2)[0] != strings.SplitN(extVersion, ".", 2)[0] {
		return errors.Errorf(versionMismatchMsg, extVersion, ext.ID, mtaVersion)
	}

	return nil
}

// GetExtensionFilePaths returns the MTA extension descriptor full paths
func (ep *Loc) GetExtensionFilePaths() []string {
	paths := make([]string, len(ep.ExtensionFileNames))
	for i, fileName := range ep.ExtensionFileNames {
		paths[i] = ep.GetMtaExtYamlPath(fileName)
	}
	return paths
}

// ParseExtFile returns a reference to the MTA extension descriptor object of the extension file.
func (ep *Loc) ParseExtFile(extFileName string) (*mta.EXT, error) {
	yamlContent, err := ReadExt(ep, extFileName)
	if err != nil {
		return nil, err
	}

	// Parse MTA extension file
	mtaExt, err := mta.UnmarshalExt(yamlContent)
	if err != nil {
		return mtaExt, errors.Wrapf(err, parseExtFileFailed, ep.GetMtaExtYamlPath(extFileName))
	}
	return mtaExt, nil
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
