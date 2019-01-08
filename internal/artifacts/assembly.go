package artifacts

import (
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/contenttype"
	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/mta"
)

const (
	moduleEntry    = "MTA-Module"
	requiredEntry  = "MTA-Requires"
	resourceEntry  = "MTA-Resource"
	dirContentType = "text/directory"
)

// Assembly - assembles mta project into .mtar
func Assembly(source, target string, wdGetter func() (string, error)) error {

	logs.Logger.Info("assembly started")

	// initialize location
	loc, err := dir.Location(source, target, "dep", wdGetter)
	if err != nil {
		return errors.Wrap(err, "assembly failed when initializing location")
	}

	// create temporary folder
	tmpDir := loc.GetTargetTmpDir()
	err = os.Mkdir(tmpDir, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "assembly failed when creating the %s folder", tmpDir)
	}

	logs.Logger.Infof("the %s temporary folder created", tmpDir)

	// get mta object
	mta, err := loc.ParseFile()
	if err != nil {
		return errors.Wrap(err, "assembly failed when parsing .mtad file")
	}

	// get entries to be assembled
	entries, err := getAssembledEntries(loc, mta)
	if err != nil {
		return errors.Wrap(err, "assembly failed when getting entries")
	}

	// copy assembled entries to the temporary folder
	err = dir.CopyEntries(getEntriesInfo(entries), loc.SourcePath, tmpDir)
	if err != nil {
		return errors.Wrapf(err, "assembly failed when copying entries from the %s folder to the %s folder",
			loc.SourcePath, tmpDir)
	}
	logs.Logger.Infof("%v files/folders copied", len(entries))

	// generate the .mtad file
	err = genMtad(mta, loc, true, "")
	if err != nil {
		return err
	}
	logs.Logger.Info("the .mtad file copied into META-INF folder of the temporary folder")

	// generate the manifest file
	err = genManifest(loc.GetManifestPath(), entries)
	if err != nil {
		return errors.Wrap(err, "assembly failed when generating the manifest file")
	}

	logs.Logger.Info("the manifest file created in META-INF folder of the temporary folder")

	// archive the temporary folder into the .mtar file
	err = dir.Archive(tmpDir, filepath.Join(target, mta.ID+mtarSuffix))
	if err != nil {
		return errors.Wrap(err, "assembly failed when archiving")
	}
	logs.Logger.Info("the .mtar file created from the temporary folder")

	// cleanup the temporary folder
	err = ExecuteCleanup(loc.GetSource(), loc.GetTarget(), "dep", os.Getwd)
	if err != nil {
		return errors.Wrap(err, "assembly failed when cleaning the temporary folder")
	}
	logs.Logger.Info("assembly finished successfully")
	return nil
}

func getAssembledEntries(loc dir.ISourceModule, mta *mta.MTA) ([]entry, error) {

	contentTypes, err := contenttype.GetContentTypes()
	if err != nil {
		return nil, errors.Wrap(err, "assembly failed when unmarshalling content types")
	}

	var entries []entry
	for _, m := range mta.Modules {
		if m.Path != "" {
			entry, err := getFileInfo(loc, m.Name, m.Path, moduleEntry, &contentTypes)
			if err != nil {
				return nil, getPathError(err, m.Path)
			}
			entries = append(entries, *entry)
		}
		if m.Requires != nil {
			for _, rm := range m.Requires {
				entries, err = addPathFromParameters(loc, rm.Name, requiredEntry, entries, &contentTypes, rm.Parameters)
				if err != nil {
					return entries, err
				}
			}
		}
	}
	for _, r := range mta.Resources {
		entries, err = addPathFromParameters(loc, r.Name, resourceEntry, entries, &contentTypes, r.Parameters)
		if err != nil {
			return entries, err
		}
	}
	return entries, nil
}

func getEntriesInfo(entries []entry) []os.FileInfo {
	var filesInfos []os.FileInfo
	for _, e := range entries {
		filesInfos = append(filesInfos, *e.file)
	}
	return filesInfos
}

func getPathError(err error, path string) error {
	return errors.Wrapf(err, "assembly failed when searching the %s path", path)
}

func addPathFromParameters(loc dir.ISourceModule, name, entryType string, entries []entry, contentTypes *contenttype.ContentTypes,
	params map[string]interface{}) ([]entry, error) {

	if params != nil {
		if entryPath, ok := params["path"]; ok {
			entry, err := getFileInfo(loc, name, entryPath.(string), entryType, contentTypes)
			if err != nil {
				return nil, getPathError(err, entryPath.(string))
			}
			return append(entries, *entry), nil
		}
	}
	return entries, nil
}

func getFileInfo(loc dir.ISourceModule, entryName, entryPath, entryType string, contentTypes *contenttype.ContentTypes) (*entry, error) {
	fullPath := loc.GetSourceModuleDir(entryPath)
	file, err := os.Stat(fullPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting information of the %s path", entryPath)
	}
	var contentType string
	if file.IsDir() {
		contentType = dirContentType
	} else {
		ext := path.Ext(entryPath)
		contentType, err = contenttype.GetContentType(contentTypes, ext)
		if err != nil {
			return nil, err
		}
	}
	return &entry{EntryName: entryName, EntryType: entryType, file: &file, ContentType: contentType, EntryPath: entryPath}, nil

}
