package metainfo

import (
	"fmt"
	"io"
	"log"
	"os"
	"io/ioutil"
	"cloud-mta-build-tool/cmd/mta/models"
	"cloud-mta-build-tool/cmd/constants"
	"cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/mta"
	"cloud-mta-build-tool/cmd/platform"
	"cloud-mta-build-tool/cmd/mta/converter"
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
	MetaInf         = "/META-INF"
	Manifest        = "MANIFEST.MF"
	Mtad            = "mtad.yaml"
	NewLine         = "\n"
	ContentType     = "Content-Type: "
	MtaModule       = "MTA-Module: "
	ModuleName      = "Name: "
	ApplicationZip  = "application/zip"
	ManifestVersion = "Manifest-Version: 1.0"
)

//Set the MANIFEST.MF file
func setManifetDesc(file io.Writer, mtaStr []*models.Modules, modules []string) {
	// TODO create dynamically
	fmt.Fprint(file, ManifestVersion+NewLine)
	// TODO set the version from external config for automatic version bump during release
	fmt.Fprint(file, "Created-By: SAP Application Archive Builder 0.0.1")
	for _, mod := range mtaStr {
		//Print only the required module to support the partial build
		if len(modules) > 0 && mod.Name == modules[0] {
			printToFile(file, mod)
			break
		} else if len(modules) == 0 {
			//Print all the modules
			printToFile(file, mod)
		}
	}
}

func printToFile(file io.Writer, mtaStr *models.Modules) {
	fmt.Fprint(file, NewLine)
	fmt.Fprint(file, NewLine)
	fmt.Fprint(file, ModuleName+mtaStr.Name+constants.DataZip)
	fmt.Fprint(file, NewLine)
	fmt.Fprint(file, MtaModule+mtaStr.Name)
	fmt.Fprint(file, NewLine)
	fmt.Fprint(file, ContentType+ApplicationZip)
}

func GenMetaInf(tmpDir string, mtaStr models.MTA, modules []string) {
	// Create META-INF folder under the mtar folder
	dir.CreateDirIfNotExist(tmpDir + MetaInf)
	//Load platform configuration file

	platformCfg := platform.Parse(platform.PlatformConfig)
	// Modify MTAD object according to platform types
	//Todo platform should provided as command parameter
	converter.ConvertTypes(mtaStr, platformCfg, "cf")
	// Create readable Yaml before writing to file
	mtad := mta.Marshal(mtaStr)
	// Write back the MTAD to the META-INF folder
	err := ioutil.WriteFile(tmpDir+MetaInf+constants.PathSep+Mtad, mtad, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	// Create MANIFEST.MF file
	file := dir.CreateFile(tmpDir + MetaInf + constants.PathSep + Manifest)
	// Set the MANIFEST.MF file
	setManifetDesc(file, mtaStr.Modules, modules)
}
