package gen

import (
	"testing"
	fs "cloud-mta-build-tool/cmd/fsys"
	"os"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/cmd/constants"
)

func basicMakeAndValidate(t *testing.T, path, yamlFilename, makeFilename, expectedMakeFileExtension string) {
	err := makeFile(path, yamlFilename, path, makeFilename, "make_verbose.txt")
	makeFullName := path + constants.PathSep + makeFilename
	if err != nil {
		os.Remove(makeFullName)
		t.Error(err)
	}
	actual, err := ioutil.ReadFile(makeFullName + expectedMakeFileExtension)
	assert.Nil(t, err)
	expected, _ := ioutil.ReadFile(path + constants.PathSep + "ExpectedMakeFile")
	assert.Equal(t, expected, actual)
}

func removeMakefile(t *testing.T, path, makeFilename string) {
	err := os.Remove(path + constants.PathSep + makeFilename)
	assert.Nil(t, err)
}

func TestMake(t *testing.T) {

	path := fs.GetPath() + constants.PathSep + "testdata"
	makeFilename := "MakeFileTest"

	logs.Logger = logs.NewLogger()

	type testInfo = struct {
		name         string
		filename     string
		testExecutor func(t *testing.T, path, yamlFilename, makeFilename string)
	}

	for _, testInfo := range []testInfo{
		{"SanityTest", "mta.yaml",
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				basicMakeAndValidate(t, path, yamlFilename, makeFilename, "")
				removeMakefile(t, path, makeFilename)
			}},
		{"Yaml file not exists", "YamlNotExists",
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				err := makeFile(path, yamlFilename, path, makeFilename, "make.txt")
				assert.NotNil(t, err)
			}},
		{"Yaml file exists but not answers YAML format", "ExpectedMakeFile",
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				err := makeFile(path, yamlFilename, path, makeFilename, "make.txt")
				assert.NotNil(t, err)
			}},
		{"Make runs twice, 2 files created - with and without extension", "mta.yaml",
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				basicMakeAndValidate(t, path, yamlFilename, makeFilename, "")
				basicMakeAndValidate(t, path, yamlFilename, makeFilename, ".mta")
				removeMakefile(t, path, makeFilename)
				removeMakefile(t, path, makeFilename+".mta")
			}},
		{"public Make testing", "mta.yaml",
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				err := Make()
				assert.NotNil(t, err)
			}},
		{"Template is wrong", "mta.yaml",
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				err := makeFile(path, yamlFilename, path, makeFilename, "testdata"+constants.PathSep+"WrongMakeTmpl.txt")
				assert.NotNil(t, err)
			}},
		{"Template is empty", "mta.yaml",
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				err := makeFile(path, yamlFilename, path, makeFilename, "testdata"+constants.PathSep+"emptyMakeTmpl.txt")
				removeMakefile(t, path, makeFilename)
				assert.NotNil(t, err)
			}},
	} {
		t.Run(testInfo.name, func(t *testing.T) {
			testInfo.testExecutor(t, path, testInfo.filename, makeFilename)
		})
	}
}
