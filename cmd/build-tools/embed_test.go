package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"gotest.tools/assert"
)

func Test_main(t *testing.T) {
	os.Args = []string{"app", "-source=./testdata/cfg.yaml", "-target=./testdata/cfg.go", "-package=testpackage", "-name=Config"}
	main()
	actualContent, _ := ioutil.ReadFile("./testdata/cfg.go")
	expectedContent, _ := ioutil.ReadFile("./testdata/goldenCfg.go")
	assert.Equal(t, strings.Replace(string(expectedContent), "0xd, ", "", -1), strings.Replace(string(actualContent), "0xd, ", "", -1))
	os.RemoveAll("./testdata/cfg.go")
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
