package main

import (
	"os"

	cmd "github.com/SAP/cloud-mta-build-tool/cmd"
)

func main() {
	// Execute CLI Root commands
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
