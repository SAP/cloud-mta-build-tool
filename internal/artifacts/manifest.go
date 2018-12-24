package artifacts

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cloud-mta-build-tool/internal/version"
	"cloud-mta-build-tool/mta"

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

// setManifestDesc - Set the MANIFEST.MF file
func setManifestDesc(file io.Writer, mtaStr []*mta.Module, modules []string) error {
	// TODO create dynamically
	_, err := fmt.Fprint(file, manifestVersion+newLine)
	if err != nil {
		return errors.Wrap(err, "failed to generate the manifest file when printing manifest version")
	}
	v, err := version.GetVersion()
	if err != nil {
		return errors.Wrap(err, "failed to generate the manifest file when getting the CLI version")
	}
	_, err = fmt.Fprintf(file, "Created-By: SAP Application Archive Builder %v", v.CliVersion)
	if err != nil {
		return errors.Wrap(err, "failed to generate the manifest file when printing the CLI version")
	}
	for _, mod := range mtaStr {
		if moduleDefined(mod.Name, modules) {
			err := printToFile(file, mod)
			if err != nil {
				return errors.Wrapf(err, "failed to generate the manifest file when printing module <%v>", mod.Name)
			}
		}
	}
	return nil
}

// moduleDefined - checks if module defined in the list
func moduleDefined(module string, modules []string) bool {
	if modules == nil || len(modules) == 0 {
		return true
	}
	for _, m := range modules {
		if m == module {
			return true
		}
	}
	return false
}

// printToFile - Print to manifest.mf file
func printToFile(file io.Writer, mtaStr *mta.Module) error {
	_, err := fmt.Fprint(file, newLine+newLine, filepath.ToSlash(moduleName+mtaStr.Name+dataZip),
		newLine, mtaModule+mtaStr.Name, newLine, contentType+applicationZip)
	return err
}
