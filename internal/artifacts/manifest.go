package artifacts

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/internal/tpl"
	"cloud-mta-build-tool/internal/version"
	"cloud-mta-build-tool/mta"
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
	applicationZip = "application/zip"
	pathSep        = string(os.PathSeparator)
	dataZip        = pathSep + "data.zip"
)

type entry struct {
	EntryName   string
	EntryType   string
	file        *os.FileInfo
	ContentType string
	EntryPath   string
}

// setManifestDesc - Set the MANIFEST.MF file
func setManifestDesc(ep dir.ITargetArtifacts, mtaStr []*mta.Module, modules []string) error {

	var entries []entry
	for _, mod := range mtaStr {
		if moduleDefined(mod.Name, modules) {
			moduleEntry := entry{
				EntryName:   mod.Name,
				EntryPath:   filepath.ToSlash(mod.Name + dataZip),
				ContentType: applicationZip,
				EntryType:   moduleEntry,
			}
			entries = append(entries, moduleEntry)
		}
	}
	return genManifest(ep.GetManifestPath(), entries)
}

func genManifest(manifestPath string, entries []entry) (rerr error) {

	v, err := version.GetVersion()
	if err != nil {
		return errors.Wrap(err, "failed to generate the manifest file when getting the CLI version")
	}

	funcMap := template.FuncMap{
		"Entries":    entries,
		"CliVersion": v.CliVersion,
	}
	out, err := os.Create(manifestPath)
	defer func() {
		errClose := out.Close()
		if errClose != nil && rerr == nil {
			rerr = errors.Wrap(errClose, "failed to generate the manifest file when closing the manifest file")
		}
	}()
	if err != nil {
		return errors.Wrap(err, "failed to generate the manifest file when creating the manifest file")
	}
	t := template.Must(template.New("template").Parse(string(tpl.Manifest)))
	err = t.Execute(out, funcMap)
	if err != nil {
		return errors.Wrap(err, "failed to generate the manifest file when populating the content")
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
