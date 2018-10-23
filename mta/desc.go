package mta

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"cloud-mta-build-tool/internal/fsys"
)

// The deployment descriptor shall be located within the META-INF folder of the JAR.
// The file MANIFEST.MF shall contain at least a name section for each MTA module contained in the archive.
// Following the JAR specification, the value of a name must be a relative path to a file or directory,
// or an absolute URL referencing data outside the archive.
// It is required to add a row MTA-module: <modulename> to each name section which corresponds to an MTA module,
// to bind archive file locations to module names as used in the deployment descriptor.
// The name sections with the MTA module attribute indicates the path to the file or directory which represents a module within the archive
// This used by deploy service to track the build project

const (
	metaInf         = "META-INF"
	manifest        = "MANIFEST.MF"
	mtadYaml        = "mtad.yaml"
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
func setManifetDesc(file io.Writer, mtaStr []*Modules, modules []string) {
	// TODO create dynamically
	fmt.Fprint(file, manifestVersion+newLine)
	// TODO set the version from external config for automatic version bump during release
	fmt.Fprint(file, "Created-By: SAP Application Archive Builder 0.0.1")
	for _, mod := range mtaStr {
		// Print only the required module to support the partial build
		if len(modules) > 0 && mod.Name == modules[0] {
			printToFile(file, mod)
			break
		} else if len(modules) == 0 {
			// Print all the modules
			printToFile(file, mod)
		}
	}
}


// Print to manifest.mf file
func printToFile(file io.Writer, mtaStr *Modules) {
	fmt.Fprint(file, newLine)
	fmt.Fprint(file, newLine)
	fmt.Fprint(file, filepath.ToSlash(moduleName+mtaStr.Name+dataZip))
	fmt.Fprint(file, newLine)
	fmt.Fprint(file, mtaModule+mtaStr.Name)
	fmt.Fprint(file, newLine)
	fmt.Fprint(file, contentType+applicationZip)
}

// GenMtad -Generate mtad.yaml from mta file and configuration file
func GenMtad(mtaStr MTA, targetPath string, convertTypes func(mtaStr MTA)) error {
	// Create META-INF folder under the mtar folder
	targetBasePath := dir.Path{targetPath}
	createDirIfNotExist(targetBasePath.GetFullPath(metaInf))
	convertTypes(mtaStr)
	// Create readable Yaml before writing to file
	mtad, err := Marshal(mtaStr)
	// Write back the MTAD to the META-INF folder
	if err == nil {
		err = ioutil.WriteFile(targetBasePath.GetFullPath(metaInf, mtadYaml), mtad, os.ModePerm)
	}
	return err
}


// GenMetaInfo -Generate MetaInfo (MANIFEST.MF file)
func GenMetaInfo(tmpDir string, mtaStr MTA, modules []string, convertTypes func(mtaStr MTA)) error {
	err := GenMtad(mtaStr, tmpDir, convertTypes)
	if err == nil {
		// Create MANIFEST.MF file
		file, _ := createFile(dir.Path{tmpDir}.GetFullPath(metaInf, manifest))
		defer file.Close()
		// Set the MANIFEST.MF file
		setManifetDesc(file, mtaStr.Modules, modules)
	}
	return err
}

// CreateDirIfNotExist - Create new dir
func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Failed to create dir %s ", err)
		}
	}
	return nil
}

// CreateFile - create new file
func createFile(path string) (file *os.File, err error) {
	file, err = os.Create(path) // Truncates if file already exists
	if err != nil {
		return nil, fmt.Errorf("Failed to create file %s ", err)
	}
	// /defer file.Close()
	return file, nil
}
