package tpl

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"

	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"

	"github.com/stretchr/testify/assert"
)

func basicMakeAndValidate(t *testing.T, path, yamlFilename, makeFilename, expectedMakeFilename, expectedMakeFileExtension string) {
	err := makeFile(path, makeFilename, "make_verbose.txt")
	makeFullName := path + pathSep + makeFilename
	if err != nil {
		os.Remove(makeFullName)
		t.Error(err)
	}
	actual, err := ioutil.ReadFile(makeFullName + expectedMakeFileExtension)
	assert.Nil(t, err)
	expected, _ := ioutil.ReadFile(path + pathSep + expectedMakeFilename)
	assert.Equal(t, expected, actual)
}

func removeMakefile(t *testing.T, path, makeFilename string) {
	err := os.Remove(path + pathSep + makeFilename)
	assert.Nil(t, err)
}

func TestMake(t *testing.T) {

	path := fs.GetPath() + pathSep + "testdata"
	makeFilename := "MakeFileTest"
	var expectedMakeFilename string
	switch runtime.GOOS {
	case "linux":
		expectedMakeFilename = "ExpectedMakeFileLinux"
	case "darwin":
		expectedMakeFilename = "ExpectedMakeFileMac"
	case "windows":
		expectedMakeFilename = "ExpectedMakeFileWindows"
	}

	logs.Logger = logs.NewLogger()

	type testInfo = struct {
		name         string
		filename     string
		testExecutor func(t *testing.T, path, yamlFilename, makeFilename string)
	}

	for _, ti := range []testInfo{
		{"SanityTest", "mta.yaml",
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				basicMakeAndValidate(t, path, yamlFilename, makeFilename, expectedMakeFilename, "")
				removeMakefile(t, path, makeFilename)
			}},
		{"Yaml file not exists", "YamlNotExists",
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				err := makeFile(path+"1", makeFilename, "make.txt")
				assert.NotNil(t, err)
			}},
		{"Yaml file exists but not answers YAML format", expectedMakeFilename,
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				err := makeFile(path, makeFilename, "make.txt")
				assert.NotNil(t, err)
			}},
		{"Make runs twice, 2 files created - with and without extension", "mta.yaml",
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				basicMakeAndValidate(t, path, yamlFilename, makeFilename, expectedMakeFilename, "")
				basicMakeAndValidate(t, path, yamlFilename, makeFilename, expectedMakeFilename, ".mta")
				removeMakefile(t, path, makeFilename)
				removeMakefile(t, path, makeFilename+".mta")
			}},
		{"public Make testing", "mta.yaml",
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				var flag string
				err := Make(flag)
				assert.NotNil(t, err)
			}},
		{"Template is wrong", "mta.yaml",
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				err := makeFile(path, makeFilename, "testdata"+pathSep+"WrongMakeTmpl.txt")
				assert.NotNil(t, err)
			}},
		{"Template is empty", "mta.yaml",
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				err := makeFile(path, makeFilename, "testdata"+pathSep+"emptyMakeTmpl.txt")
				removeMakefile(t, path, makeFilename)
				assert.NotNil(t, err)
			}},
	} {
		t.Run(ti.name, func(t *testing.T) {
			ti.testExecutor(t, path, ti.filename, makeFilename)
		})
	}
}

func Test_makeMode(t *testing.T) {

	type args struct {
		mode string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Default template - Generate user template according to command params ",
			args: args{
				mode: "",
			},
			want:    "make_default.txt",
			wantErr: false,
		},

		{
			name: "Verbose template - Generate user template according to command params ",
			args: args{
				mode: "verbose",
			},
			want:    "make_verbose.txt",
			wantErr: false,
		},
		{
			name: "Unsupported command - Generate user template according to command params",
			args: args{
				mode: "--test",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := makeMode(tt.args.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeMode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("makeMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringInSlice(t *testing.T) {
	type args struct {
		a    string
		list []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "find string in slice",
			args: args{
				a:    "--test1",
				list: []string{"--test1", "foo"},
			},
			want: true,
		},
		{
			name: "find string in slice",
			args: args{
				a:    "--test",
				list: []string{"--test1", "bar"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringInSlice(tt.args.a, tt.args.list); got != tt.want {
				t.Errorf("stringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
