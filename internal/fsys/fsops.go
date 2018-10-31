package dir

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	pathSep = string(os.PathSeparator)
)

// CreateDirIfNotExist - Create new dir
func CreateDirIfNotExist(dir string) error {
	var err error
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	}
	return err
}

// Archive module and mtar artifacts,
// compatible with the JAR specification
// to support the spec requirements
// Source Path to be zipped
// Target artifact
//TODO add more comments
func Archive(sourcePath, targetArchivePath string) error {

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
	defer zipfile.Close()

	// create archive writer
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

	// pack files of source into archive
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
			// care of UNIX-style separators of path in header
			header.Name = filepath.ToSlash(GetRelativePath(path, baseDir))
		}

		// compress file
		header.Method = zip.Deflate

		// add new header and file to archive
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
