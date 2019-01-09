package dir

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/logs"
)

// CreateDirIfNotExist - Create new dir
func CreateDirIfNotExist(dir string) error {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	} else if !info.IsDir() {
		err = fmt.Errorf("creation of the %v folder failed because a file exists with the same name", dir)
	}
	return err
}

// Archive module and mtar artifacts,
// compatible with the JAR specification
// to support the spec requirements
// Source Path to be zipped
// Target artifact
func Archive(sourcePath, targetArchivePath string) (e error) {

	// check that folder to be packed exist
	info, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	// create archive file
	zipfile, err := os.Create(targetArchivePath)
	if err != nil {
		return err
	}
	defer func() {
		errClose := zipfile.Close()
		if errClose != nil && e == nil {
			e = errClose
		}
	}()

	// create archive writer
	archive := zip.NewWriter(zipfile)
	defer func() {
		errClose := archive.Close()
		if errClose != nil && e == nil {
			e = errClose
		}
	}()

	// Skip headers to support jar archive structure
	var baseDir string
	if info.IsDir() {
		baseDir = sourcePath
	} else {
		baseDir = filepath.Base(sourcePath)
	}

	if baseDir != "" {
		baseDir += string(os.PathSeparator)
	}

	err = walk(sourcePath, baseDir, archive)
	if err != nil {
		return err
	}
	return err
}

func walk(sourcePath string, baseDir string, archive *zip.Writer) (e error) {
	// pack files of source into archive
	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			// care of UNIX-style separators of path in header
			header.Name = filepath.ToSlash(getRelativePath(path, baseDir))
		}

		// compress file
		header.Method = zip.Deflate

		// add new header and file to archive
		writer, err := archive.CreateHeader(header)
		if err == nil {
			file, e := os.Open(path)
			if e == nil {
				defer func() {
					errClose := file.Close()
					if errClose != nil && e == nil {
						e = errClose
					}
				}()
				_, err = io.Copy(writer, file)
				return err
			}
		}
		return err
	})
}

// CreateFile - create new file
func CreateFile(path string) (file *os.File, err error) {
	file, err = os.Create(path) // Truncates if file already exists
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("creation of the %s file failed", path))
	}
	// The caller needs to use defer.close
	return file, err
}

// CopyDir - copy directory content
func CopyDir(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("copying of the %v folder to the %v folder failed because the source is not a folder", src, dst)
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	return CopyEntries(entries, src, dst)
}

// CopyByPatterns - copy files/directories according to patterns
// from source folder to target folder
// patterns are relative to source folder
func CopyByPatterns(source, target string, patterns []string) error {

	if patterns == nil || len(patterns) == 0 {
		return nil
	}

	logs.Logger.Infof("copying the patterns [%v,...] from the %v folder to the %v folder started",
		patterns[0], source, target)

	infoTargetDir, err := os.Stat(target)
	if err != nil {
		err = os.MkdirAll(target, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err,
				"copying the patterns [%v,...] from the %v folder to the %v folder failed when creating the target folder",
				patterns[0], source, target)
		}
		logs.Logger.Infof("the %v folder has been created", target)

	} else if !infoTargetDir.IsDir() {
		return errors.Errorf(
			"copying the patterns [%v,...] from the %v folder to the %v folder failed because the target is not a folder",
			patterns[0], source, target)
	}

	for _, pattern := range patterns {
		err = copyByPattern(source, target, pattern)
		if err != nil {
			return err
		}
	}

	logs.Logger.Infof("copying the patterns [%v,...] from the %v folder to the %v folder finished successfully",
		patterns[0], source, target)
	return nil
}

// copyByPattern - copy files/directories according to pattern
func copyByPattern(source, target, pattern string) error {
	logs.Logger.Infof("copying the %v pattern from the %v folder to the %v folder started",
		pattern, source, target)
	// build full pattern concatenating source path and pattern
	fullPattern := filepath.Join(source, strings.Replace(pattern, "./", "", -1))
	// get all entries matching the pattern
	sourceEntries, err := filepath.Glob(fullPattern)
	if err != nil {
		return errors.Wrapf(err,
			"copying the %v pattern from the %v folder to the %v folder failed when getting matching entries",
			pattern, source, target)
	}

	for _, sourceEntry := range sourceEntries {
		info, err := os.Stat(sourceEntry)
		if err != nil {
			return errors.Wrapf(err,
				"copying the %v pattern from the %v folder to the %v folder failed when getting the status of the source entry: %v",
				pattern, source, target, sourceEntry)
		}
		targetEntry := filepath.Join(target, filepath.Base(sourceEntry))
		if info.IsDir() {
			err = CopyDir(sourceEntry, targetEntry)
		} else {
			err = CopyFile(sourceEntry, targetEntry)
		}
		if err != nil {
			return errors.Wrapf(err,
				"copying the %v pattern from the %v folder to the %v folder failed when copying the %v entry to the %v entry",
				pattern, source, target, sourceEntry, targetEntry)
		}
	}
	logs.Logger.Infof(
		"copying the %v pattern from the %v folder to the %v folder finished successfully",
		pattern, source, target)
	return nil
}

// CopyEntries - copies entries (files and directories) from source to destination folder
func CopyEntries(entries []os.FileInfo, src, dst string) error {

	var err error
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// execute recursively
			err = CopyDir(srcPath, dstPath)
		} else {
			// Todo check posix compatibility
			if entry.Mode()&os.ModeSymlink != 0 {
				fmt.Println(
					fmt.Sprintf(
						"copying of the entries from the %v folder to the %v folder skipped the %v entry because its mode is a symbolic link",
						src, dst, entry.Name()),
					src, dst, entry.Name())
				continue
			}

			err = CopyFile(srcPath, dstPath)
		}
		if err != nil {
			break
		}
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
		if e := in.Close(); e != nil && rerr == nil {
			rerr = e
		}
	}()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if e := out.Close(); e != nil && rerr == nil {
			rerr = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, si.Mode())
	return err
}

// getRelativePath - Remove the basePath from the fullPath and get only the relative
func getRelativePath(fullPath, basePath string) string {
	return strings.TrimPrefix(fullPath, basePath)
}

// Read returns mta byte slice.
func Read(ep IMtaYaml) ([]byte, error) {
	fileFullPath := ep.GetMtaYamlPath()
	// Read MTA file
	yamlFile, err := ioutil.ReadFile(fileFullPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read the %v file", fileFullPath)
	}
	return yamlFile, nil
}

// ReadExt returns mta extension byte slice.
func ReadExt(ep IMtaExtYaml, platform string) ([]byte, error) {
	fileFullPath := ep.GetMtaExtYamlPath(platform)
	// Read MTA extension file
	yamlFile, err := ioutil.ReadFile(fileFullPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read the %v file", fileFullPath)
	}
	return yamlFile, err
}
