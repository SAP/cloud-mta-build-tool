package artifacts

import (
	"os"
	"path/filepath"
	"text/template"

	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/internal/tpl"
	"cloud-mta-build-tool/internal/version"

	"github.com/SAP/cloud-mta/mta"

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
func setManifestDesc(ep dir.ITargetArtifacts, targetPathGetter dir.ITargetPath, mtaStr []*mta.Module, mtaResources []*mta.Resource, modules []string) error {
	var entries []entry
	for _, mod := range mtaStr {
		if moduleDefined(mod.Name, modules) {
			moduleEntry := entry{
				EntryName:   mod.Name,
				EntryPath:   getModulePath(mod, targetPathGetter),
				ContentType: getContentType(targetPathGetter, getModulePath(mod, targetPathGetter)),
				EntryType:   moduleEntry,
			}
			entries = append(entries, moduleEntry)
		}
	}

	for _, resource := range mtaResources {
		if resource.Parameters["path"] == nil {
			continue
		}
		resourceEntry := entry{
			EntryName:   resource.Name,
			EntryPath:   getResourcePath(resource),
			ContentType: getContentType(targetPathGetter, getResourcePath(resource)),
			EntryType:   resourceEntry,
		}
		entries = append(entries, resourceEntry)
	}

	for _, mod := range mtaStr {
		if moduleDefined(mod.Name, modules) {
			requiredDependenciesWithPath := getRequiredDependencies(mod)
			requiredDependencyEntries := buildEntries(targetPathGetter, mod, requiredDependenciesWithPath)
			entries = append(entries, requiredDependencyEntries...)
		}
	}

	return genManifest(ep.GetManifestPath(), entries)
}

func buildEntries(targetPathGetter dir.ITargetPath, module *mta.Module, requiredDependencies []mta.Requires) []entry {
	result := make([]entry, 0)
	for _, requiredDependency := range requiredDependencies {
		requiredDependencyEntry := entry{
			EntryName:   module.Name + "/" + requiredDependency.Name,
			EntryPath:   requiredDependency.Parameters["path"].(string),
			ContentType: getContentType(targetPathGetter, requiredDependency.Parameters["path"].(string)),
			EntryType:   requiredEntry,
		}
		result = append(result, requiredDependencyEntry)
	}
	return result
}

func getContentType(targetPathGetter dir.ITargetPath, path string) string {
	if targetPathGetter == nil {
		return applicationZip
	}
	info, err := os.Stat(filepath.Join(targetPathGetter.GetTargetTmpDir(), path))
	if err != nil {
		return ""
	}

	if info.IsDir() {
		return dirContentType
	}

	return applicationZip
}

func getRequiredDependencies(module *mta.Module) []mta.Requires {
	result := make([]mta.Requires, 0)
	for _, requiredDependency := range module.Requires {
		if requiredDependency.Parameters["path"] != nil {
			result = append(result, requiredDependency)
		}
	}
	return result
}

func getResourcePath(resource *mta.Resource) string {
	return resource.Parameters["path"].(string)
}

func getModulePath(module *mta.Module, targetPathGetter dir.ITargetPath) string {
	if targetPathGetter == nil {
		return filepath.ToSlash(module.Name + dataZip)
	}
	loc := targetPathGetter.(*dir.Loc)
	if existsModuleZipInDirectories(module, []string{loc.GetSource(), loc.GetTargetTmpDir()}) {
		return filepath.ToSlash(module.Name + dataZip)
	}
	return module.Path
}

func existsModuleZipInDirectories(module *mta.Module, directories []string) bool {
	for _, directory := range directories {
		if _, err := os.Stat(filepath.Join(directory, filepath.ToSlash(module.Name+dataZip))); !os.IsNotExist(err) {
			return true
		}
	}
	return false
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
