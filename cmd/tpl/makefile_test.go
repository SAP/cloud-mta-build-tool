package tpl

import (
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
	"testing"

	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"

	"github.com/stretchr/testify/assert"
)

func basicMakeAndValidate(t *testing.T, path, yamlFilename, makeFilename, expectedMakeFilename, expectedMakeFileExtension string) {
	tpl := tplCfg{tplName: "make_verbose.txt", relPath: "/testdata", pre: basePreVerbose, post: basePostVerbose}
	err := makeFile(makeFilename, tpl)
	makeFullName := path + pathSep + makeFilename
	if err != nil {
		os.Remove(makeFullName)
		t.Error(err)
	}
	actual, err := ioutil.ReadFile(makeFullName + expectedMakeFileExtension)
	assert.Nil(t, err)
	expected, _ := ioutil.ReadFile(path + pathSep + expectedMakeFilename)
	assert.Equal(t, removeSpecialSymbols(expected), removeSpecialSymbols(actual))
}

func removeSpecialSymbols(b []byte) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	s := string(b)
	s = strings.Replace(s, "0xd, ", "", -1)
	s = reg.ReplaceAllString(s, "")
	return s
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
				tpl := tplCfg{tplName: "make.txt"}
				err := makeFile(makeFilename, tpl)
				assert.NotNil(t, err)
			}},
		{"Yaml file exists but not answers YAML format", expectedMakeFilename,
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				tpl := tplCfg{tplName: "make.txt"}
				err := makeFile(makeFilename, tpl)
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
				err := Make("errorMode")
				assert.NotNil(t, err)
			}},
		{"Template is wrong", "mta.yaml",
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				tpl := tplCfg{tplName: "testdata" + pathSep + "WrongMakeTmpl.txt"}
				err := makeFile(makeFilename, tpl)
				assert.NotNil(t, err)
			}},
		{"Template is empty", "mta.yaml",
			func(t *testing.T, path, yamlFilename, makeFilename string) {
				tpl := tplCfg{tplName: "testdata" + pathSep + "emptyMakeTmpl.txt"}
				err := makeFile(makeFilename, tpl)
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
		name string
		args args

		want    tplCfg
		wantErr bool
	}{
		{
			name: "Default template - Generate user template according to command params ",
			args: args{
				mode: "",
			},
			want:    tplCfg{tplName: makeDefaultTpl, pre: basePreDefault, post: basePostDefault},
			wantErr: false,
		},

		{
			name: "Verbose template - Generate user template according to command params ",
			args: args{
				mode: "verbose",
			},
			want:    tplCfg{tplName: makeVerboseTpl, pre: basePreVerbose, post: basePostVerbose},
			wantErr: false,
		},
		{
			name: "Unsupported command - Generate user template according to command params",
			args: args{
				mode: "--test",
			},
			want:    tplCfg{},
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
