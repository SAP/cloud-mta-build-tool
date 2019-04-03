package artifacts

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta/mta"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/conttype"
	"github.com/SAP/cloud-mta-build-tool/internal/tpl"
	"github.com/SAP/cloud-mta-build-tool/internal/version"
	"github.com/SAP/cloud-mta-build-tool/internal/buildops"
	"github.com/SAP/cloud-mta-build-tool/internal/commands"
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
	pathSep        = string(os.PathSeparator)
	dataZip        = pathSep + "data.zip"
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
func setManifestDesc(ep dir.ITargetArtifacts, targetPathGetter dir.ITargetPath, mtaStr []*mta.Module,
	mtaResources []*mta.Resource, modules []string, onlyModules bool) error {

	contentTypes, err := conttype.GetContentTypes()
	if err != nil {
		return errors.Wrap(err,
			"failed to generate the manifest file when getting the content types from the configuration")
	}

	var entries []entry
	for _, mod := range mtaStr {
		if !moduleDefined(mod.Name, modules) || mod.Name == "" {
			continue
		}
		_, defaultBuildResult, err := commands.CommandProvider(*mod)
		if err != nil {
			return err
		}
		modulePath, err := getModulePath(mod, targetPathGetter, defaultBuildResult)
		if err != nil {
			return err
		}
		contentType, err := getContentType(targetPathGetter, modulePath, contentTypes)
		if err != nil {
			return errors.Wrapf(err,
				`failed to generate the manifest file when getting the "%s" module content type`, mod.Name)
		}

		entries = addModuleEntry(entries, mod, contentType, modulePath)

		if onlyModules {
			continue
		}
		requiredDependenciesWithPath := getRequiredDependencies(mod)
		requiredDependencyEntries, err :=
			buildEntries(targetPathGetter, mod, requiredDependenciesWithPath, contentTypes)
		if err != nil {
			return errors.Wrapf(err,
				`failed to generate the manifest file when building the required entries of the "%s" module`,
				mod.Name)
		}
		entries = append(entries, requiredDependencyEntries...)
	}

	if !onlyModules {
		resourcesEntries, err := getResourcesEntries(targetPathGetter, mtaResources, contentTypes)
		if err != nil {
			return err
		}
		entries = append(entries, resourcesEntries...)
	}

	return genManifest(ep.GetManifestPath(), entries)
}

func addModuleEntry(entries []entry, module *mta.Module, contentType, modulePath string) []entry {
	result := entries

	if modulePath != "" {
		moduleEntry := entry{
			EntryName:   module.Name,
			EntryPath:   modulePath,
			ContentType: contentType,
			EntryType:   moduleEntry,
		}
		result = append(entries, moduleEntry)
	}
	return result
}

func getResourcesEntries(targetPathGetter dir.ITargetPath, resources []*mta.Resource,
	contentTypes *conttype.ContentTypes) ([]entry, error) {
	var entries []entry
	for _, resource := range resources {
		if resource.Name == "" || resource.Parameters["path"] == nil {
			continue
		}
		contentType, err := getContentType(targetPathGetter, getResourcePath(resource), contentTypes)
		if err != nil {
			return nil, errors.Wrapf(err,
				`failed to generate the manifest file when getting the "%s" resource content type`, resource.Name)
		}
		resourceEntry := entry{
			EntryName:   resource.Name,
			EntryPath:   getResourcePath(resource),
			ContentType: contentType,
			EntryType:   resourceEntry,
		}
		entries = append(entries, resourceEntry)
	}
	return entries, nil
}

func buildEntries(targetPathGetter dir.ITargetPath, module *mta.Module,
	requiredDependencies []mta.Requires, contentTypes *conttype.ContentTypes) ([]entry, error) {
	result := make([]entry, 0)
	for _, requiredDependency := range requiredDependencies {
		contentType, err :=
			getContentType(targetPathGetter, requiredDependency.Parameters["path"].(string), contentTypes)
		if err != nil {
			return nil, err
		}
		requiredDependencyEntry := entry{
			EntryName:   module.Name + "/" + requiredDependency.Name,
			EntryPath:   requiredDependency.Parameters["path"].(string),
			ContentType: contentType,
			EntryType:   requiredEntry,
		}
		result = append(result, requiredDependencyEntry)
	}
	return result, nil
}

func getContentType(targetPathGetter dir.ITargetPath, path string, contentTypes *conttype.ContentTypes) (string, error) {
	fullPath := filepath.Join(targetPathGetter.GetTargetTmpDir(), path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return "", fmt.Errorf(`the "%s" path does not exist; the content type was not defined`, fullPath)
	}

	if info.IsDir() {
		return dirContentType, nil
	}

	extension := filepath.Ext(fullPath)
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
	return resource.Parameters["path"].(string)
}

func getModulePath(module *mta.Module, targetPathGetter dir.ITargetPath, defaultBuildResult string) (string, error) {
	loc := targetPathGetter.(*dir.Loc)

	// get build results path - defined in build-params property or in 
	buildResultPath, buildResultDefined, err := buildops.GetBuildResultsPath(loc, module, defaultBuildResult)
	if err != nil {
		return "", err
	}
	if buildResultPath == "" {
		// module path not defined
		return module.Path, nil
	} else if buildResultDefined {
		return filepath.Join(module.Name, filepath.Base(buildResultPath)), nil
	} else if existsModuleZipInDirectories(module, []string{loc.GetSource(), loc.GetTargetTmpDir()}) {
		return filepath.ToSlash(module.Name + dataZip), nil
	}
	return filepath.Base(buildResultPath), nil
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
		rerr = dir.CloseFile(out, rerr)
	}()
	if err != nil {
		return errors.Wrap(err, "failed to generate the manifest file when initializing it")
	}
	return populateManifest(out, funcMap)
}

func populateManifest(file io.Writer, funcMap template.FuncMap) error {
	t := template.Must(template.New("template").Parse(string(tpl.Manifest)))
	err := t.Execute(file, funcMap)
	if err != nil {
		return errors.Wrap(err, "failed to generate the manifest file when populating the content")
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
