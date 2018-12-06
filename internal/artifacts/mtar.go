package artifacts

import (
	"path/filepath"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"

	"github.com/pkg/errors"
)

const (
	mtarSuffix = ".mtar"
)

// GenerateMtar - generate mtar archive from the build artifacts
func GenerateMtar(ep *dir.Loc) error {
	logs.Logger.Info("MTAR Generation started")
	// get MTA object
	m, err := dir.ParseFile(ep)
	if err != nil {
		return errors.Wrap(err, "MTAR Generation failed on MTA parsing")
	}
	// get target temporary folder to be archived
	targetTmpDir, err := ep.GetTargetTmpDir()
	if err != nil {
		return errors.Wrap(err, "MTAR Generation failed on getting target temp directory")
	}
	// get target directory - where mtar will be saved
	targetDir, err := ep.GetTarget()
	if err != nil {
		return errors.Wrap(err, "MTAR Generation failed on getting target directory")
	}
	// archive building artifacts to mtar
	err = dir.Archive(targetTmpDir, filepath.Join(targetDir, m.ID+mtarSuffix))
	if err != nil {
		return errors.Wrap(err, "MTAR Generation failed on MTAR archiving")
	}
	logs.Logger.Info("MTAR Generation successfully finished")
	return nil
}
