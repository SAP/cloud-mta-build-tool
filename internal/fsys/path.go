package dir

import (
	"os"
	"path/filepath"
	"strings"
)

// Path - path to a files
type Path struct {
	Path string
}

// GetCurrentPath - get current Path
func GetCurrentPath() (string, error) {
	// TODO should get also from user
	return os.Getwd()
}

// GetFullPath - get full Path (currentPath + relPath)
func GetFullPath(relPath ...string) (string, error) {
	path, err := GetCurrentPath()
	if err == nil {
		pathElements := []string{path}
		path = filepath.Join(append(pathElements, relPath...)...)
	}
	return path, err
}

// GetFullPath - relative to the basePath
func (basePath Path) GetFullPath(relPath ...string) string {
	path := basePath.Path
	pathElements := []string{path}
	path = filepath.Join(append(pathElements, relPath...)...)
	return path
}

// GetArtifactsPath - the Path where all the build file will be saved
func GetArtifactsPath(path string) (string, error) {
	currentPath, err := GetFullPath(path)
	var artifactsPath string
	if err == nil {
		_, file := filepath.Split(currentPath)
		// append to the currentPath the file name
		artifactsPath = filepath.Join(currentPath, file)
	}
	return artifactsPath, err
}

// GetRelativePath - remove the basePath from the fullPath and get only the relative
func GetRelativePath(fullPath, basePath string) string {
	return strings.TrimPrefix(fullPath, basePath)
}
