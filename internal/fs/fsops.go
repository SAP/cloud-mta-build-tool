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
		err = fmt.Errorf("CreateDirIfNotExist failed - %v file exists", dir)
	}
	return err
}

// Archive module and mtar artifacts,
// compatible with the JAR specification
// to support the spec requirements
// Source Path to be zipped
// Target artifact
func Archive(sourcePath, targetArchivePath string) (rerr error) {

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
		if errClose != nil {
			rerr = errClose
		}
	}()

	// create archive writer
	archive := zip.NewWriter(zipfile)
	defer func() {
		errClose := archive.Close()
		if errClose != nil {
			rerr = errClose
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
		return errors.Wrap(err, "Archiving error")
	}
	return err
}

func walk(sourcePath string, baseDir string, archive *zip.Writer) (rerr error) {
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
					if errClose != nil {
						rerr = errClose
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
		return nil, fmt.Errorf("Failed to create file %s ", err)
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
		return fmt.Errorf("The provided source %s is not a directory ", src)
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

	return copyEntries(entries, src, dst)
}

// CopyByPatterns - copy files/directories according to patterns
// from source folder to target folder
// patterns are relative to source folder
func CopyByPatterns(source, target string, patterns []string) error {

	if patterns == nil || len(patterns) == 0 {
		return nil
	}

	logs.Logger.Infof("Copy by patterns started. Source <%v> target <%v>", source, target)

	infoTargetDir, err := os.Stat(target)
	if err != nil {
		err = os.MkdirAll(target, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "Copy by patterns [%v,...] failed on creating directory %v", patterns[0], target)
		}
		logs.Logger.Infof("Directory <%v> created", target)

	} else if !infoTargetDir.IsDir() {
		return errors.Errorf("Copy by patterns [%v,...] failed. Target-path %v is not a folder", patterns[0], target)
	}

	for _, pattern := range patterns {
		err = copyByPattern(source, target, pattern)
		if err != nil {
			return err
		}
	}

	logs.Logger.Info("Copy by patterns successfully finished.")
	return nil
}

// copyByPattern - copy files/directories according to pattern
func copyByPattern(source, target, pattern string) error {
	logs.Logger.Infof("Copy by pattern <%v> started.", pattern)
	// build full pattern concatenating source path and pattern
	fullPattern := filepath.Join(source, strings.Replace(pattern, "./", "", -1))
	// get all entries matching the pattern
	sourceEntries, err := filepath.Glob(fullPattern)
	if err != nil {
		return errors.Wrapf(err, "Copy by pattern %v failed on getting matching entries", pattern)
	}

	for _, sourceEntry := range sourceEntries {
		info, err := os.Stat(sourceEntry)
		if err != nil {
			return errors.Wrapf(err, "Copy by pattern %v failed on getting status of source entry %v", pattern, sourceEntry)
		}
		targetEntry := filepath.Join(target, filepath.Base(sourceEntry))
		if info.IsDir() {
			err = CopyDir(sourceEntry, targetEntry)
		} else {
			err = CopyFile(sourceEntry, targetEntry)
		}
		if err != nil {
			return errors.Wrapf(err, "Copy by pattern %v failed on copy of %v to %v", pattern, sourceEntry, targetEntry)
		}
	}
	logs.Logger.Infof("Copy by pattern <%v> successfully finished.", pattern)
	return nil
}

// copyEntries - copies entries (files and directories) from source to destination folder
func copyEntries(entries []os.FileInfo, src, dst string) error {

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
				fmt.Println("MBT: SymbolicLink ignored")
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
		if e := in.Close(); e != nil {
			rerr = e
		}
	}()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if e := out.Close(); e != nil {
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
		return nil, errors.Wrap(err, "Error reading the MTA file")
	}
	return yamlFile, nil
}

// ReadExt returns mta extension byte slice.
func ReadExt(ep IMtaExtYaml, platform string) ([]byte, error) {
	fileFullPath := ep.GetMtaExtYamlPath(platform)
	// Read MTA extension file
	yamlFile, err := ioutil.ReadFile(fileFullPath)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading the EXT file")
	}
	return yamlFile, err
}
