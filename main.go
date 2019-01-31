package main

import (
	"github.com/SAP/cloud-mta-build-tool/cmd"
	"os"
)

func main() {
	// Execute CLI Root commands
	err := commands.Execute()
	if err != nil {
		os.Exit(1)
	}
}
