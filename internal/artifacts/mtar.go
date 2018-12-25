package artifacts

import (
	"path/filepath"

	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/internal/logs"
)

const (
	mtarSuffix = ".mtar"
)

// ExecuteGenMtar - generates MTAR
func ExecuteGenMtar(source, target, desc string, wdGetter func() (string, error)) error {
	logs.Logger.Info("generation of .mtar started")
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrap(err, "generation of .mtar failed when initializing location")
	}
	err = generateMtar(loc, loc)
	if err != nil {
		return err
	}
	logs.Logger.Info("generation of .mtar finished successfully")
	return nil
}

// generateMtar - generate mtar archive from the build artifacts
func generateMtar(targetLoc dir.ITargetPath, parser dir.IMtaParser) error {
	// get MTA object
	m, err := parser.ParseFile()
	if err != nil {
		return errors.Wrap(err, "generation of the .mtar file failed when parsing the mta file")
	}
	// get target temporary folder to be archived
	targetTmpDir := targetLoc.GetTargetTmpDir()
	// get target directory - where mtar will be saved
	targetDir := targetLoc.GetTarget()
	// archive building artifacts to mtar
	err = dir.Archive(targetTmpDir, filepath.Join(targetDir, m.ID+mtarSuffix))
	if err != nil {
		return errors.Wrap(err, "generation of .mtar failed when archiving")
	}
	return nil
}
