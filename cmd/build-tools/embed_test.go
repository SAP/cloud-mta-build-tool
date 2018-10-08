package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"gotest.tools/assert"
)

func Test_main(t *testing.T) {
	os.Args = []string{"app", "-source=./testdata/cfg.yaml", "-target=./testdata/cfg.go", "-package=testpackage", "-name=Config"}
	main()
	actualContent, _ := ioutil.ReadFile("./testdata/cfg.go")
	expectedContent, _ := ioutil.ReadFile("./testdata/goldenCfg.go")
	assert.Equal(t, removeSpecialSymbols(expectedContent), removeSpecialSymbols(actualContent))
	os.RemoveAll("./testdata/cfg.go")
}

func removeSpecialSymbols(b []byte) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(string(b), "")
}

func Test_mainNegative(t *testing.T) {
	os.Args = []string{"app", "-source=./testdata/cfgNotExisting.yaml", "-target=./testdata/cfg.go", "-package=testpackage", "-name=Config"}
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	main()
}
