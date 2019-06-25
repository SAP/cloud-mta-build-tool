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
	dataZip        = "data.zip"
	moduleEntry    = "MTA-Module"
	requiredEntry  = "MTA-Requires"
	resourceEntry  = "MTA-Resource"
	dirContentType = "text/directory"

	// ModuleMsgProperty - module message property
	ModuleMsgProperty = "module"
	// ResourceMsgProperty - resource message property
	ResourceMsgProperty = "resource"

	// UnknownContentTypeMsg - unknown content type message
	UnknownContentTypeMsg = `failed to generate the manifest file when getting the "%s" %s content type`
)

type entry struct {
	EntryName   string
	EntryType   string
	ContentType string
	EntryPath   string
}

// setManifestDesc - Set the MANIFEST.MF file
func setManifestDesc(source dir.ISourceModule, ep dir.ITargetArtifacts, targetPathGetter dir.ITargetPath, depDesc bool, mtaStr []*mta.Module,
	mtaResources []*mta.Resource, modules []string) error {

	contentTypes, err := conttype.GetContentTypes()
	if err != nil {
		return errors.Wrap(err,
			"failed to generate the manifest file when getting the content types from the configuration")
	}

	entries, err := getModulesEntries(source, targetPathGetter, depDesc, mtaStr, contentTypes, modules)
	if err != nil {
		return err
	}

	resourcesEntries, err := getResourcesEntries(targetPathGetter, mtaResources, contentTypes)
	if err != nil {
		return err
	}
	entries = append(entries, resourcesEntries...)

	return genManifest(ep.GetManifestPath(), entries)
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
	contentTypes *conttype.ContentTypes, modules []string) ([]entry, error) {

	var entries []entry
	for _, mod := range moduleList {
		if !moduleDefined(mod.Name, modules) || mod.Name == "" {
			continue
		}
		_, defaultBuildResult, err := commands.CommandProvider(*mod)
		if err != nil {
			return nil, err
		}
		modulePath, _, err := buildops.GetModuleTargetArtifactPath(source, targetPathGetter, depDesc, mod, defaultBuildResult, true)
		if err != nil {
			return nil, errors.Wrapf(err,
				`failed to generate the manifest file when getting the artifact path of the "%s" module`, mod.Name)
		}

		if modulePath != "" {
			contentType, err := getContentType(modulePath, contentTypes)
			if err != nil {
				return nil, errors.Wrapf(err, UnknownContentTypeMsg, mod.Name, ModuleMsgProperty)
			}

			// get relative path of the module entry (excluding leading slash)
			moduleEntryPath := strings.Replace(modulePath, targetPathGetter.GetTargetTmpDir(), "", -1)[1:]
			entries = addModuleEntry(entries, mod, contentType, moduleEntryPath)
		}

		requiredDependenciesWithPath := getRequiredDependencies(mod)
		requiredDependencyEntries, err := buildEntries(targetPathGetter, mod, requiredDependenciesWithPath, contentTypes)
		if err != nil {
			return nil, errors.Wrapf(err,
				`failed to generate the manifest file when building the required entries of the "%s" module`,
				mod.Name)
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
			return nil, errors.Wrapf(err, UnknownContentTypeMsg, resource.Name, ResourceMsgProperty)
		}
		resourceEntry := entry{
			EntryName:   resource.Name,
			EntryPath:   resourceRelativePath,
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
		contentType, err := getContentType(filepath.Join(target.GetTargetTmpDir(), requiredDependency.Parameters["path"].(string)), contentTypes)
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

func getContentType(path string, contentTypes *conttype.ContentTypes) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf(`the "%s" path does not exist; the content type was not defined`, path)
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
