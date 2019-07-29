package dir

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

// CreateDirIfNotExist - Create new dir
func CreateDirIfNotExist(dir string) error {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	} else if err == nil && !info.IsDir() {
		err = errors.Errorf(FolderCreationFailedMsg, dir)
	}
	return err
}

// Archive module and mtar artifacts,
// compatible with the JAR specification
// to support the spec requirements
// Source Path to be zipped
// Target artifact
func Archive(sourcePath, targetArchivePath string, ignore []string) (e error) {

	// check that folder to be packed exist
	info, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	// create folder of archive file if not exists
	err = CreateDirIfNotExist(filepath.Dir(targetArchivePath))
	if err != nil {
		return errors.Wrapf(err, archivingFailedOnCreateFolderMsg, filepath.Dir(targetArchivePath))
	}

	// create archive file
	zipfile, err := os.Create(targetArchivePath)
	if err != nil {
		return err
	}
	defer func() {
		e = CloseFile(zipfile, e)
	}()

	// create archive writer
	archive := zip.NewWriter(zipfile)
	defer func() {
		e = CloseFile(archive, e)
	}()

	// Skip headers to support jar archive structure
	var baseDir string
	if info.IsDir() {
		baseDir = sourcePath
	} else {
		baseDir = filepath.Dir(sourcePath)
	}

	if !strings.HasSuffix(baseDir, string(os.PathSeparator)) {
		baseDir += string(os.PathSeparator)
	}

	ignoreMap, err := getIgnoredEntries(ignore, sourcePath)
	if err != nil {
		return err
	}

	err = walk(sourcePath, baseDir, archive, ignoreMap)
	return err
}

// getIgnoresMap - getIgnores Helper
func getIgnoredEntries(ignore []string, sourcePath string) (map[string]interface{}, error) {
	ignoredEntriesMap := map[string]interface{}{}
	for _, ign := range ignore {
		path := filepath.Join(sourcePath, ign)
		entries, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			ignoredEntriesMap[entry] = nil
		}
	}
	return ignoredEntriesMap, nil
}

// CloseFile - closes file
// error handling takes into account error of the calling function
func CloseFile(file io.Closer, err error) error {
	errClose := file.Close()
	if errClose != nil && err == nil {
		return errClose
	}
	return err
}

func walk(sourcePath string, baseDir string, archive *zip.Writer, ignore map[string]interface{}) error {

	// pack files of source into archive
	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) (e error) {
		if err != nil {
			return err
		}

		if _, ok := ignore[path]; ok {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return
		}

		// Don't add the base folder to the zip
		if info.IsDir() {
			pathWithSlash := path
			if !strings.HasSuffix(pathWithSlash, string(os.PathSeparator)) {
				pathWithSlash += string(os.PathSeparator)
			}
			if pathWithSlash == baseDir {
				return
			}
		}

		// Path in zip should be with slashes (in all operating systems)
		pathInZip := filepath.ToSlash(getRelativePath(path, baseDir))

		// Folders must end with "/"
		if info.IsDir() {
			pathInZip += "/"
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			header.Method = zip.Deflate
		}
		header.Name = pathInZip

		// add new header and file to archive
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			e = CloseFile(file, e)
		}()

		_, err = io.Copy(writer, file)

		return err
	})
}

// CreateFile - create new file
func CreateFile(path string) (file *os.File, err error) {
	file, err = os.Create(path) // Truncates if file already exists
	if err != nil {
		return nil, errors.Wrapf(err, fileCreationFailedMsg, path)
	}
	// The caller needs to use defer.close
	return file, err
}

// CopyDir - copy directory content
func CopyDir(src string, dst string, withParents bool, copyDirEntries func(entries []os.FileInfo, src, dst string) error) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	_, err := os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if !withParents && err != nil {
		err = os.Mkdir(dst, os.ModePerm)
	} else if err != nil {
		err = CreateDirIfNotExist(dst)
	}
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	return copyDirEntries(entries, src, dst)
}

// CopyByPatterns - copy files/directories according to patterns
// from source folder to target folder
// patterns are relative to source folder
func CopyByPatterns(source, target string, patterns []string) error {

	if len(patterns) == 0 {
		return nil
	}

	logs.Logger.Infof(copyByPatternMsg, patterns[0], source, target)

	infoTargetDir, err := os.Stat(target)
	if err != nil {
		err = CreateDirIfNotExist(target)
		if err != nil {
			return errors.Wrapf(err, copyByPatternFailedOnCreateMsg, patterns[0], source, target, target)
		}
		logs.Logger.Infof(folderCreatedMsg, target)

	} else if !infoTargetDir.IsDir() {
		return errors.Errorf(copyByPatternFailedOnTargetMsg, patterns[0], source, target, target)
	}

	for _, pattern := range patterns {
		err = copyByPattern(source, target, pattern)
		if err != nil {
			return err
		}
	}

	return nil
}

// copyByPattern - copy files/directories according to pattern
func copyByPattern(source, target, pattern string) error {
	logs.Logger.Infof(`copying the "%s" pattern from the "%s" folder to the "%s" folder`,
		pattern, source, target)
	// build full pattern concatenating source path and pattern
	fullPattern := filepath.Join(source, strings.Replace(pattern, "./", "", -1))
	// get all entries matching the pattern
	sourceEntries, err := filepath.Glob(fullPattern)
	if err != nil {
		return errors.Wrapf(err, copyByPatternFailedOnMatchMsg, pattern, source, target, pattern)
	}

	err = copyEntries(sourceEntries, source, target, pattern)
	if err != nil {
		return err
	}

	return nil
}

func copyEntries(entries []string, source, target, pattern string) error {
	for _, entry := range entries {
		info, err := os.Stat(entry)
		if err != nil {
			return errors.Wrapf(err, copyFailedOnGetStatusMsg, pattern, source, target, entry)
		}
		targetEntry := filepath.Join(target, filepath.Base(entry))
		if info.IsDir() {
			err = CopyDir(entry, targetEntry, true, CopyEntries)
		} else {
			err = CopyFileWithMode(entry, targetEntry, info.Mode())
		}
		if err != nil {
			return errors.Wrapf(err, copyFailedMsg, pattern, source, target, entry, targetEntry)
		}
	}
	return nil
}

// CopyEntries - copies entries (files and directories) from source to destination folder
func CopyEntries(entries []os.FileInfo, src, dst string) error {

	if len(entries) == 0 {
		return nil
	}
	for _, entry := range entries {
		var err error
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// execute recursively
			err = CopyDir(srcPath, dstPath, false, CopyEntries)
		} else {
			// Todo check posix compatibility
			if entry.Mode()&os.ModeSymlink != 0 {
				logs.Logger.Infof(skipSymbolicLinkMsg, src, dst, entry.Name())
			} else {
				err = CopyFileWithMode(srcPath, dstPath, entry.Mode())
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// CopyEntriesInParallel - copies entries (files and directories) from source to destination folder in parallel
func CopyEntriesInParallel(entries []os.FileInfo, src, dst string) (rerr error) {

	// limit parallel processes
	const maxOpenFiles = 5

	if len(entries) == 0 {
		return nil
	}
	// handle parallel processes with limited slice of semaphores
	sem := make(chan bool, maxOpenFiles)
	for _, entry := range entries {
		// if copy failed stop processing
		if rerr != nil {
			break
		}
		sem <- true
		go func(e os.FileInfo) {

			// free place in semaphores at the end of routine
			defer func() { <-sem }()

			var err error
			srcPath := filepath.Join(src, e.Name())
			dstPath := filepath.Join(dst, e.Name())

			if e.IsDir() {
				// execute recursively
				err = CopyDir(srcPath, dstPath, false, CopyEntriesInParallel)
			} else {
				// Todo check posix compatibility
				if e.Mode()&os.ModeSymlink != 0 {
					logs.Logger.Infof(skipSymbolicLinkMsg, src, dst, e.Name())
				} else {
					err = CopyFileWithMode(srcPath, dstPath, e.Mode())
				}
			}
			if err != nil {
				rerr = err
			}
		}(entry)
	}
	// wait for the end of all running go routines
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	return
}

// CopyFileWithMode - copy file content using file mode parameter
func CopyFileWithMode(src, dst string, mode os.FileMode) (rerr error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		rerr = CloseFile(in, rerr)
	}()

	out, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer func() {
		rerr = CloseFile(out, rerr)
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return err

}

// CopyFile - copy file content
func CopyFile(src, dst string) (rerr error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		rerr = CloseFile(in, rerr)
	}()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		rerr = CloseFile(out, rerr)
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	err = changeTargetMode(src, dst)

	return err

}

func changeTargetMode(source, target string) error {
	si, err := os.Stat(source)
	if err != nil {
		return err
	}
	return os.Chmod(target, si.Mode())
}

// getRelativePath - Remove the basePath from the fullPath and get only the relative
func getRelativePath(fullPath, basePath string) string {
	if basePath == "" {
		return fullPath
	}
	return strings.TrimPrefix(fullPath, basePath)
}

// Read returns mta byte slice.
func Read(ep IMtaYaml) ([]byte, error) {
	// Read MTA file
	return readFile(ep.GetMtaYamlPath())
}

func readFile(file string) ([]byte, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrapf(err, ReadFailedMsg, file)
	}
	s := string(content)
	s = strings.Replace(s, "\r\n", "\r", -1)
	content = []byte(s)
	return content, nil
}

// ReadExt returns mta extension byte slice.
func ReadExt(ep IMtaExtYaml, platform string) ([]byte, error) {
	return readFile(ep.GetMtaExtYamlPath(platform))
}
