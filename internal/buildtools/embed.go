package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/pkg/errors"
)

type configInfo struct {
	PackageName string
	VarName     string
	Data        string
}

var templatePath = "./internal/buildtools"

// This code is executed during the go:generate command
// Reading the config files and generate byte array
func main() {
	// define commands for execution
	inFile := flag.String("source", "", "source")
	outFile := flag.String("target", "", "target")
	pkg := flag.String("package", "main", "package name")
	name := flag.String("name", "File", "identifier to use for the embedded data")
	flag.Parse()
	err := genConf(*inFile, *outFile, *pkg, *name)
	handleError(err)
}

func genConf(source string, target, packageName, varName string) (e error) {
	// Read the config file
	inData, err := ioutil.ReadFile(source)
	if err != nil {
		return errors.Wrapf(err, readFailedMsg, source)
	}
	out, err := os.Create(target)
	defer func() {
		e = dir.CloseFile(out, e)
	}()
	if err != nil {
		return errors.Wrapf(err, createFailedMsg, target)
	}
	t := template.Must(template.New("config.tpl").ParseFiles(filepath.Join(templatePath, "config.tpl")))
	err = t.Execute(out, configInfo{PackageName: packageName, VarName: varName, Data: fmt.Sprintf("%#v", inData)})
	if err != nil {
		return errors.Wrapf(err, populationFailedMsg)
	}
	return nil
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
}
