package commands

import (
	"errors"
	"fmt"
	"path/filepath"

	fs "cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"

	"cloud-mta-build-tool/internal/platform"
	"cloud-mta-build-tool/mta"
)

// generate build metadata artifacts
func generateMeta(relPath string, args []string) error {
	return processMta("Metadata creation", relPath, args, func(file []byte, args []string) error {
		// Parse MTA file
		m, err := mta.ParseToMta(file)
		if err == nil {
			// Generate meta info dir with required content
			err = mta.GenMetaInfo(args[0], *m, args[1:], func(mtaStr mta.MTA) {
				convertTypes(mtaStr)
			})
		}
		return err
	})
}

// generate mtar archive from the build artifacts
func generateMtar(relPath string, args []string) error {
	return processMta("MTAR generation", relPath, args, func(file []byte, args []string) error {
		// Create MTAR from the building artifacts
		m, err := mta.ParseToMta(file)
		if err == nil {
			err = fs.Archive(filepath.Join(args[0]), filepath.Join(args[1], m.Id+mtarSuffix))
		}
		return err
	})
}

// convert types to appropriate target platform types
func convertTypes(mtaStr mta.MTA) {
	// Load platform configuration file
	platformCfg := platform.Parse(platform.PlatformConfig)
	// Modify MTAD object according to platform types
	// Todo platform should provided as command parameter
	platform.ConvertTypes(mtaStr, platformCfg, "cf")
}

// process mta.yaml file
func processMta(processName string, relPath string, args []string, process func(file []byte, args []string) error) error {
	logs.Logger.Info("Starting " + processName)
	s := &mta.Source{Path: relPath, Filename: "mta.yaml"}
	mf, err := s.Readfile()
	if err == nil {
		err = process(mf, args)
		if err == nil {
			logs.Logger.Info(processName + " finish successfully ")
		}
	} else {
		err = errors.New(fmt.Sprintf("MTA file not found: %s", err))
	}
	return err
}
