package dir

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"cloud-mta-build-tool/cmd/logs"
)

const (
	pathSep = string(os.PathSeparator)
)

// CreateDirIfNotExist - Create new dir
func CreateDirIfNotExist(dir string) error {
	var err error
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			logs.Logger.Error(err)
		}
	}
	return err
}

// Archive module and mtar artifacts,
// compatible with the JAR specification
// to support the spec requirements
// Source Path to be zipped
// Target artifact
func Archive(sourcePath, targetArchivePath string) error {

	info, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	zipfile, err := os.Create(targetArchivePath)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// Skip headers to support jar archive structure
	var baseDir string
	if info.IsDir() {
		baseDir = sourcePath
	} else {
		baseDir = filepath.Base(sourcePath)
	}

	if baseDir != "" {
		baseDir += pathSep
	}

	filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
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
			header.Name = GetRelativePath(path, baseDir)
			header.Name = ConvertPathToUnixFormat(header.Name)
		}

		header.Method = zip.Deflate

		writer, err := archive.CreateHeader(header)
		if err == nil {
			file, err := os.Open(path)
			if err == nil {
				defer file.Close()
				_, err = io.Copy(writer, file)
			}
		}
		return err
	})

	return err
}

// CreateFile - create new file
func CreateFile(path string) (file *os.File, err error) {
	file, err = os.Create(path) // Truncates if file already exists
	if err != nil {
		return nil, fmt.Errorf("Failed to create file %s ", err)
	}
	// /defer file.Close()
	return file, err
}

// CopyDir - copy directory content
func CopyDir(src string, dst string) (err error) {
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

func copyEntries(entries []os.FileInfo, src, dst string) error {

	var err error
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// execute recursively
			err = CopyDir(srcPath, dstPath)
		} else {
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = copyFile(srcPath, dstPath)
		}
		if err != nil {
			break
		}
	}
	return err
}

// CopyFile - copy file content
func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	err = out.Sync()
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

// DefaultTempDirFunc - Currently the generated temp dir is one lvl up from the running project
// e.g. proj -> go/src/mta tmpdir-> go/src
// The tmp dir should be deleted after the process finished or failed
// TODO delete tmp in case of failure
func DefaultTempDirFunc(path string) string {
	tmpDir, _ := ioutil.TempDir(path, "BUILD_MTAR_TEMP")
	return tmpDir
}

// Load - load the mta.yaml file
func Load(path string) (content []byte, err error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Logger.Errorf("File not found for Path %s, error is: #%v ", path, err)
		// YAML descriptor file not found abort the process
		return yamlFile, err
	}
	logs.Logger.Debugf("The file loaded successfully:" + string(yamlFile))
	return yamlFile, err
}
