package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

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
		template := "package %s\n\nvar %s []byte = %#v\n"
		s := fmt.Sprintf(template, *pkg, *name, inData)
		err = ioutil.WriteFile(*outFile, []byte(s), 0644)
	}
	if err != nil {
		panic(err.Error())
	}
}
