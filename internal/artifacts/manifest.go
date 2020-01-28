package artifacts

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta/mta"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/buildops"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/conttype"
	"github.com/SAP/cloud-mta-build-tool/internal/tpl"
	"github.com/SAP/cloud-mta-build-tool/internal/version"
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
	moduleEntry    = "MTA-Module"
	requiredEntry  = "MTA-Requires"
	resourceEntry  = "MTA-Resource"
	dirContentType = "text/directory"
)

type entry struct {
	EntryName   string
	EntryType   string
	ContentType string
	EntryPath   string
}

// setManifestDesc - Set the MANIFEST.MF file
func setManifestDesc(source dir.ISourceModule, ep dir.ITargetArtifacts, targetPathGetter dir.ITargetPath, depDesc bool, mtaStr []*mta.Module,
	mtaResources []*mta.Resource) error {

	contentTypes, err := conttype.GetContentTypes()
	if err != nil {
		return errors.Wrap(err, contentTypeCfgMsg)
	}

	entries, err := getModulesEntries(source, targetPathGetter, depDesc, mtaStr, contentTypes)
	if err != nil {
		return err
	}

	resourcesEntries, err := getResourcesEntries(targetPathGetter, mtaResources, contentTypes)
	if err != nil {
		return err
	}
	entries = append(entries, resourcesEntries...)

	// Module entries that point to the same path should be merged in the manifest
	entries = mergeDuplicateEntries(entries)

	return genManifest(ep.GetManifestPath(), entries)
}

func mergeDuplicateEntries(entries []entry) []entry {
	// Several MTA-Module entries can point to the same path. In that case, their names should in the same entry, comma-separated.
	mergedEntries := make([]entry, 0)
	modules := make(map[string]entry)
	// To keep a consistent sort order for the map entries we must keep another data structure (slice of keys by order of addition here)
	pathsOrder := make([]string, 0)

	// Add module entries to modules. Add non-module entries to mergedEntries.
	for index, entry := range entries {
		if entry.EntryType == moduleEntry {
			if existing, ok := modules[entry.EntryPath]; ok {
				existing.EntryName += ", " + entry.EntryName
				modules[entry.EntryPath] = existing
			} else {
				modules[entry.EntryPath] = entries[index]
				pathsOrder = append(pathsOrder, entry.EntryPath)
			}
		} else {
			mergedEntries = append(mergedEntries, entry)
		}
	}

	// Sort module entries by order of insertion
	moduleEntries := make([]entry, 0)
	for _, path := range pathsOrder {
		moduleEntries = append(moduleEntries, modules[path])
	}

	// Add the module entries first to the merged entries
	mergedEntries = append(moduleEntries, mergedEntries...)
	return mergedEntries
}

func addModuleEntry(entries []entry, module *mta.Module, contentType, modulePath string) []entry {
	result := entries

	if modulePath != "" {
		moduleEntry := entry{
			EntryName:   module.Name,
			EntryPath:   filepath.ToSlash(modulePath),
			ContentType: contentType,
			EntryType:   moduleEntry,
		}
		result = append(entries, moduleEntry)
	}
	return result
}

func getModulesEntries(source dir.ISourceModule, targetPathGetter dir.ITargetPath, depDesc bool, moduleList []*mta.Module,
	contentTypes *conttype.ContentTypes) ([]entry, error) {

	var entries []entry
	for _, mod := range moduleList {
		if !buildops.IfNoSource(mod) {
			_, defaultBuildResult, err := commands.CommandProvider(*mod)
			if err != nil {
				return nil, err
			}
			modulePath, _, err := buildops.GetModuleTargetArtifactPath(source, targetPathGetter, depDesc, mod, defaultBuildResult)
			if modulePath != "" && err == nil {
				_, err = os.Stat(modulePath)
			}

			if err != nil {
				return nil, errors.Wrapf(err, wrongArtifactPathMsg, mod.Name)
			}

			if modulePath != "" {
				contentType, err1 := getContentType(modulePath, contentTypes)
				if err1 != nil {
					return nil, errors.Wrapf(err1, unknownModuleContentTypeMsg, mod.Name)
				}

				// get relative path of the module entry (excluding leading slash)
				moduleEntryPath := strings.Replace(modulePath, targetPathGetter.GetTargetTmpDir(), "", 1)[1:]
				entries = addModuleEntry(entries, mod, contentType, moduleEntryPath)
			}
		}

		requiredDependenciesWithPath := getRequiredDependencies(mod)
		requiredDependencyEntries, err := buildEntries(targetPathGetter, mod, requiredDependenciesWithPath, contentTypes)
		if err != nil {
			return nil, errors.Wrapf(err, requiredEntriesProblemMsg, mod.Name)
		}
		entries = append(entries, requiredDependencyEntries...)
	}
	return entries, nil
}

func getResourcesEntries(target dir.ITargetPath, resources []*mta.Resource, contentTypes *conttype.ContentTypes) ([]entry, error) {
	var entries []entry
	for _, resource := range resources {
		if resource.Name == "" || resource.Parameters["path"] == nil {
			continue
		}
		resourceRelativePath := getResourcePath(resource)
		contentType, err := getContentType(filepath.Join(target.GetTargetTmpDir(), resourceRelativePath), contentTypes)
		if err != nil {
			return nil, errors.Wrapf(err, unknownResourceContentTypeMsg, resource.Name)
		}
		resourceEntry := entry{
			EntryName:   resource.Name,
			EntryPath:   filepath.ToSlash(resourceRelativePath),
			ContentType: contentType,
			EntryType:   resourceEntry,
		}
		entries = append(entries, resourceEntry)
	}
	return entries, nil
}

func buildEntries(target dir.ITargetPath, module *mta.Module, requiredDependencies []mta.Requires, contentTypes *conttype.ContentTypes) ([]entry, error) {
	result := make([]entry, 0)
	for _, requiredDependency := range requiredDependencies {
		depPath := requiredDependency.Parameters["path"].(string)
		contentType, err := getContentType(filepath.Join(target.GetTargetTmpDir(), depPath), contentTypes)
		if err != nil {
			return nil, err
		}
		requiredDependencyEntry := entry{
			EntryName:   module.Name + "/" + requiredDependency.Name,
			EntryPath:   filepath.ToSlash(filepath.Clean(depPath)),
			ContentType: contentType,
			EntryType:   requiredEntry,
		}
		result = append(result, requiredDependencyEntry)
	}
	return result, nil
}

func getContentType(path string, contentTypes *conttype.ContentTypes) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf(contentTypeDefMsg, path)
	}

	if info.IsDir() {
		return dirContentType, nil
	}

	extension := filepath.Ext(path)
	return conttype.GetContentType(contentTypes, extension)
}

func getRequiredDependencies(module *mta.Module) []mta.Requires {
	result := make([]mta.Requires, 0)
	for _, requiredDependency := range module.Requires {
		if requiredDependency.Parameters["path"] != nil && requiredDependency.Name != "" {
			result = append(result, requiredDependency)
		}
	}
	return result
}

func getResourcePath(resource *mta.Resource) string {
	return filepath.Clean(resource.Parameters["path"].(string))
}

func genManifest(manifestPath string, entries []entry) (rerr error) {

	v, err := version.GetVersion()
	if err != nil {
		return errors.Wrap(err, cliVersionMsg)
	}

	funcMap := template.FuncMap{
		"Entries":    entries,
		"CliVersion": v.CliVersion,
	}
	out, err := os.Create(manifestPath)
	defer func() {
		rerr = dir.CloseFile(out, rerr)
	}()
	if err != nil {
		return errors.Wrap(err, initMsg)
	}
	return populateManifest(out, funcMap)
}

func populateManifest(file io.Writer, funcMap template.FuncMap) error {
	t := template.Must(template.New("template").Parse(string(tpl.Manifest)))
	err := t.Execute(file, funcMap)
	if err != nil {
		return errors.Wrap(err, populationMsg)
	}

	return nil
}

// moduleDefined - checks if module defined in the list
func moduleDefined(module string, modules []string) bool {
	if len(modules) == 0 {
		return true
	}
	for _, m := range modules {
		if m == module {
			return true
		}
	}
	return false
}
