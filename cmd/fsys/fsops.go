package dir

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/cmd/constants"
)

// CreateDirIfNotExist - Create new dir
func CreateDirIfNotExist(dir string) string {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			logs.Logger.Error(err)
		}
	}
	return dir
}

// Archive module and mtar artifacts,
// compatible with the JAR specification
// to support the spec requirements
// Source path to zip -> params[0])
// Target artifact  -> ,params[1])
// Target path to zip -> params[2])
func Archive(params ...string) error {

	zipfile, err := os.Create(params[1])
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(params[0])
	if err != nil {
		return err
	}

	// Skip headers to support jar archive structure
	var baseDir string
	if info.IsDir(); len(params) > 2 {
		baseDir = params[2]
	} else {
		baseDir = filepath.Base(params[0])
	}

	if baseDir != "" {
		baseDir += "/"
	}

	filepath.Walk(params[0], func(path string, info os.FileInfo, err error) error {
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
			header.Name = filepath.Join(strings.TrimPrefix(path, baseDir))
		}

		header.Method = zip.Deflate

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

// CreateFile - create new file
func CreateFile(path string) *os.File {
	file, err := os.Create(path) // Truncates if file already exists
	if err != nil {
		logs.Logger.Fatalf("Failed to create file: %s , %s", path, err)
	}
	// /defer file.Close()

	if err != nil {
		logs.Logger.Fatalf("Failed to write to file: %s , %s", path, err)
	}
	return file
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
		logs.Logger.Println("The provided source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// execute recursively
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = copyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

// CopyFile - copy file content
func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// DefaultTempDirFunc - Currently the generated temp dir is one lvl up from the running project
// e.g. proj -> go/src/mta tmpdir-> go/src
// The tmp dir should be deleted after the process finished or failed
// TODO delete tmp in case of failure
func DefaultTempDirFunc(path string) string {
	tmpDir, _ := ioutil.TempDir(path, constants.TempFolder)
	return tmpDir
}
