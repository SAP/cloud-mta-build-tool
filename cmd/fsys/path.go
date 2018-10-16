package dir

import (
	"os"
	"path/filepath"
	"strings"
)

type Path struct {
	Path string
}

// GetCurrentPath - get current Path
func GetCurrentPath() (string, error) {
	// TODO should get also from user
	return os.Getwd()
}

func GetFullPath(relPath ...string) (string, error) {
	path, err := GetCurrentPath()
	if err == nil {
		pathElements := []string{path}
		path = filepath.Join(append(pathElements, relPath...)...)
	}
	return path, err
}

func (basePath Path) GetFullPath(relPath ...string) string {
	path := basePath.Path
	pathElements := []string{path}
	path = filepath.Join(append(pathElements, relPath...)...)
	return path
}

func GetArtifactsPath() (string, error) {
	currentPath, err := GetCurrentPath()
	var artifactsPath string
	if err == nil {
		_, file := filepath.Split(currentPath)
		artifactsPath = filepath.Join(currentPath, file)
	}
	return artifactsPath, err
}

func GetRelativePath(fullPath, basePath string) string {
	return strings.TrimPrefix(fullPath, basePath)
}

func ConvertPathToUnixFormat(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}
