package mta

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// The deployment descriptor should be located within the META-INF folder of the JAR.
// The MANIFEST.MF file should contain at least a name section for each MTA module contained in the archive.
// Following the JAR specification, the value of a name must be a relative path to a file or directory,
// or an absolute URL referencing data outside the archive.
// It is required to add a row MTA-module: <modulename> to each name section that corresponds to an MTA module, and
// to bind archive file locations to module names as used in the deployment descriptor.
// The name sections with the MTA module attribute indicate the path to the file or directory which represents a module within the archive
// This is used by the deploy service to track the build project.

const (
	newLine         = "\n"
	contentType     = "Content-Type: "
	mtaModule       = "MTA-Module: "
	moduleName      = "Name: "
	applicationZip  = "application/zip"
	manifestVersion = "manifest-Version: 1.0"
	pathSep         = string(os.PathSeparator)
	dataZip         = pathSep + "data.zip"
)

// setManifetDesc - Set the MANIFEST.MF file
func setManifetDesc(file io.Writer, mtaStr []*Modules, modules []string) error {
	// TODO create dynamically
	_, err := fmt.Fprint(file, manifestVersion+newLine)
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed")
	}
	// TODO set the version from external config for automatic version bump during release
	_, err = fmt.Fprint(file, "Created-By: SAP Application Archive Builder 0.0.1")
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed")
	}
	for _, mod := range mtaStr {
		// Print only the required module to support the partial build
		if len(modules) > 0 && mod.Name == modules[0] {
			err := printToFile(file, mod)
			if err != nil {
				return errors.Wrap(err, "Error while printing values to mtad file")
			}
			break
		} else if len(modules) == 0 {
			// Print all the modules
			err := printToFile(file, mod)
			if err != nil {
				return errors.Wrap(err, "Error while printing values to mtad file")
			}
		}
	}
	return nil
}

// Print to manifest.mf file
func printToFile(file io.Writer, mtaStr *Modules) error {
	if _, err := fmt.Fprint(file, newLine+newLine, filepath.ToSlash(moduleName+mtaStr.Name+dataZip),
		newLine, mtaModule+mtaStr.Name, newLine, contentType+applicationZip); err != nil {
		return err
	}
	return nil
}

// GenMtad generates an mtad.yaml file from a mta.yaml file and a platform configuration file.
func GenMtad(mtaStr *MTA, ep *Loc, convertTypes func(mtaStr *MTA)) error {
	// Create META-INF folder under the mtar folder
	metaPath, err := ep.GetMetaPath()
	if err != nil {
		return errors.Wrap(err, "mtad.yaml generation failed")
	}
	err = createDirIfNotExist(metaPath)
	if err != nil {
		return errors.Wrap(err, "mtad.yaml generation failed, not able to create dir")
	}
	if !ep.IsDeploymentDescriptor() {
		convertTypes(mtaStr)
	}
	// Create readable Yaml before writing to file
	mtad, err := Marshal(mtaStr)
	if err != nil {
		return errors.Wrap(err, "mtad.yaml generation failed")
	}
	mtadPath, err := ep.GetMtadPath()
	if err == nil {
		// Write back the MTAD to the META-INF folder
		err = ioutil.WriteFile(mtadPath, mtad, os.ModePerm)
	}
	if err != nil {
		return errors.Wrap(err, "mtad.yaml generation failed")
	}
	return nil
}

// GenMetaInfo generates a MANIFEST.MF file and updates the build artifacts paths for deployment purposes.
func GenMetaInfo(ep *Loc, mtaStr *MTA, modules []string, convertTypes func(mtaStr *MTA)) error {
	err := GenMtad(mtaStr, ep, convertTypes)
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed")
	}
	// Create MANIFEST.MF file
	manifestPath, err := ep.GetManifestPath()
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed")
	}
	file, err := createFile(manifestPath)
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed")
	}
	defer file.Close()
	// Set the MANIFEST.MF file
	err = setManifetDesc(file, mtaStr.Modules, modules)
	if err != nil {
		return errors.Wrap(err, "META INFO generation failed")
	}

	return nil
}

// createDirIfNotExist - Create newGn dir
func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Failed to create dir %s ", err)
		}
	}
	return nil
}

// CreateFile - create newGn file
func createFile(path string) (file *os.File, err error) {
	file, err = os.Create(path) // Truncates if file already exists
	if err != nil {
		return nil, fmt.Errorf("Failed to create file %s ", err)
	}
	return file, nil
}
