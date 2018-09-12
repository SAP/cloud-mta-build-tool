package main

import (
	"io/ioutil"
	"os"
	"testing"

	"gotest.tools/assert"
)

func Test_main(t *testing.T) {
	os.Args = []string{"app", "-source=./testdata/cfg.yaml", "-target=./testdata/cfg.go", "-package=testpackage", "-name=Config"}
	main()
	actualContent, _ := ioutil.ReadFile("./testdata/cfg.go")
	expectedContent, _ := ioutil.ReadFile("./testdata/goldenCfg.go")
	assert.Equal(t, string(expectedContent), string(actualContent))
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
