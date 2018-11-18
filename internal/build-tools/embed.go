package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	//"os"
	"text/template"
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
	pkg := flag.String("package", "main", "Package name")
	name := flag.String("name", "File", "Identifier to use for the embedded data")
	flag.Parse()
	// Read the config file
	inData, err := ioutil.ReadFile(*inFile)

	if err == nil {
		out, err := os.Create(*outFile)
		handleError(err)
		t := template.Must(template.New("config.tpl").ParseFiles(filepath.Join(templatePath, "config.tpl")))
		err = t.Execute(out, configInfo{PackageName: *pkg, VarName: *name, Data: fmt.Sprintf("%#v", inData)})
		handleError(err)
	} else {
		handleError(err)
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println("ERROR")
		panic(err.Error())
	}
}
