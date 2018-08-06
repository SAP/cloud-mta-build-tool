package metainfo

import (
	"io/ioutil"
	"log"
	"mbtv2/cmd/mta/converter"
	"mbtv2/cmd/mta/models"
	"os"

	"mbtv2/cmd/constants"
	"mbtv2/cmd/fsys"
	"mbtv2/cmd/mta"
	"io"
	"fmt"
)

// The deployment descriptor shall be located within the META-INF folder of the JAR.
// The file MANIFEST.MF shall contain at least a name section for each MTA module contained in the archive.
// Following the JAR specification, the value of a name must be a relative path to a file or directory,
// or an absolute URL referencing data outside the archive.It
// is required to add a row MTA-module: <modulename> to each name section which corresponds to an MTA module,
// to bind archive file locations to module names as used in the deployment descriptor.
// The name sections with the MTA module attribute indicates the path to the file or directory which represents a module within the archive

const (
	MetaInf         = "/META-INF"
	Manifest        = "MANIFEST.MF"
	Mtad            = "mtad.yaml"
	NewLine         = "\n"
	ContentType     = "Content-Type: application/zip"
	MtaModule       = "MTA-Module: "
	ModuleName      = "Name: "
	ManifestVersion = "Manifest-Version: 1.0 \n"
)

func setManifetDesc(file io.Writer, mtaStr []*models.Modules) {
	// TODO create dynamically
	fmt.Fprint(file, ManifestVersion)
	// TODO set the version from external config for automatic version bump during release
	fmt.Fprint(file, "Created-By: SAP Application Archive Builder 0.0.1")
	for _, mod := range mtaStr {

		fmt.Fprint(file, NewLine)
		fmt.Fprint(file, NewLine)
		fmt.Fprint(file, ModuleName+mod.Name+constants.DataZip)
		fmt.Fprint(file, NewLine)
		fmt.Fprint(file, MtaModule+mod.Name)
		fmt.Fprint(file, NewLine)
		fmt.Fprint(file, ContentType)

	}
}

func GenMetaInf(tmpDir string, mtaStr models.MTA) {
	// Create META-INF folder under the mtar folder
	dir.CreateDirIfNotExist(tmpDir + MetaInf)
	// Modify MTAD object
	converter.ModifyMtad(mtaStr)
	// Create readable Yaml before writing to file
	mtad := mta.Marshal(mtaStr)
	// Write back the mtad to the META-INF folder
	err := ioutil.WriteFile(tmpDir+MetaInf+constants.PathSep+Mtad, mtad, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	// Create MANIFEST.MF file
	file := dir.CreateFile(tmpDir + MetaInf + constants.PathSep + Manifest)
	// Set the MANIFEST.MF file
	setManifetDesc(file, mtaStr.Modules)
}
