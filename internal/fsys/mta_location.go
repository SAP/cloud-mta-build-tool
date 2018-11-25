package dir

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	dep  = "dep"
	mtad = "mtad.yaml"
)

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

// OsGetWd - get working dir
var OsGetWd = func() (string, error) {
	return os.Getwd()
}

// GetWorkingDirectory assignment
var GetWorkingDirectory = OsGetWd

// GetSource gets the processed project path;
// if it is not provided, use the current directory.
func (ep *Loc) GetSource() (string, error) {
	if ep.SourcePath == "" {
		wd, err := GetWorkingDirectory()
		if err != nil {
			return "", errors.Wrap(err, "GetSource failed")
		}
		return wd, nil
	}
	return ep.SourcePath, nil
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

// getMtaYamlFilename - Gets the MTA .yaml file name.
func (ep *Loc) getMtaYamlFilename() string {
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
	return filepath.Join(source, ep.getMtaYamlFilename()), nil
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
