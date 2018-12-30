package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

type configInfo struct {
	PackageName string
	VarName     string
	Data        string
}

var templatePath = "./internal/build-tools"

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

func genConf(source string, target, packageName, varName string) error {
	// Read the config file
	inData, err := ioutil.ReadFile(source)
	if err != nil {
		return errors.Wrapf(err, "configuration generation failed when reading the %s file", source)
	}
	out, err := os.Create(target)
	if err != nil {
		return errors.Wrapf(err, "configuration generation failed when creating the %s file", target)
	}
	t := template.Must(template.New("config.tpl").ParseFiles(filepath.Join(templatePath, "config.tpl")))
	err = t.Execute(out, configInfo{PackageName: packageName, VarName: varName, Data: fmt.Sprintf("%#v", inData)})
	errClose := out.Close()
	if err != nil {
		if errClose != nil {
			return errors.Wrapf(err, "configuration generation failed; failed to close the %s file bacause: %s",
				target, errClose.Error())
		}
		return errors.Wrap(err, "configuration generation failed")
	} else if errClose != nil {
		return errors.Wrapf(err, "configuration generation failed when closing the %s file", target)
	}
	return nil
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
}
