package dir

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPath(t *testing.T) {
  
	t.Parallel()
	path, err := GetCurrentPath()
	assert.NotEmpty(t, path)
	assert.Nil(t, err)
}

func TestGetFullPath(t *testing.T) {
	currentPath,_:= os.Getwd()
	basePath := Path{currentPath}
	got := basePath.GetFullPath(filepath.Join("testdata", "mtahtml5"))
	expected := filepath.Join(basePath.Path, "testdata", "mtahtml5")
	if got != expected {
		t.Errorf("expected output %v, actual %v", expected, got)
	}
}

func TestPath_GetFullPath(t *testing.T) {
	currentPath,_:= os.Getwd()
	tests := []struct {
		input        []string
		expected     string
	}{
		{
			input: []string{"testdata"},
			expected: filepath.Join(currentPath, "testdata"),
		},
		{
			input: []string{"testdata", "mtahtml5"},
			expected: filepath.Join(currentPath, "testdata", "mtahtml5"),
		},
		{
			input: []string{"testdata", "level2"},
			expected: filepath.Join(currentPath, "testdata", "level2"),
		},
	}
	for _, tt := range tests {
		got := getFullPath(tt.input...)
		if got != tt.expected {
			t.Errorf("expected output %v, actual %v", tt.expected, got)
		}
	}
}

func TestGetArtifactsPath(t *testing.T) {
	currentPath,_ := os.Getwd()
	got,_ := GetArtifactsPath()
	expected := filepath.Join(currentPath, "fsys")
	if got != expected {
		t.Errorf("expected output %v, actual %v", expected, got)
	}
}

func TestGetRelativePath(t *testing.T) {
	tests := []struct {
		fullPath     string
		basePath     string
		expected     string
	}{
		{
			fullPath: filepath.Join("https:"," ", "github.com", "SAP", "cloud-mta-build-tool"),
			basePath: filepath.Join("https:"," ", "github.com"),
			expected: filepath.Join(" ", "SAP", "cloud-mta-build-tool"),
		},
		{
			fullPath: filepath.Join("https:"," ", "github.com", "SAP", "cloud-mta-build-tool"),
			basePath: filepath.Join("https:"," ", "github.com", "SAP"),
			expected: filepath.Join(" ", "cloud-mta-build-tool"),
		},
	}
	for _, tt := range tests {
		got := GetRelativePath(tt.fullPath, tt.basePath)
		assert.NotEqual(t, tt.expected, got)
	}
}