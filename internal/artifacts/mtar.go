package artifacts

import (
	"path/filepath"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/mta"

	"github.com/pkg/errors"
)

const (
	mtarSuffix = ".mtar"
)

// GenerateMtar - generate mtar archive from the build artifacts
func GenerateMtar(ep *dir.Loc) error {
	logs.Logger.Info("MTAR Generation started")
	err := processMta("MTAR generation", ep, []string{}, func(file []byte, args []string) error {
		// read MTA
		m, err := mta.Unmarshal(file)
		if err != nil {
			return errors.Wrap(err, "MTA Process failed on yaml parsing")
		}
		targetTmpDir, err := ep.GetTargetTmpDir()
		if err != nil {
			return errors.Wrap(err, "MTA Process failed on getting target temp directory")
		}
		targetDir, err := ep.GetTarget()
		if err != nil {
			return errors.Wrap(err, "MTA Process failed on getting target directory")
		}
		// archive building artifacts to mtar
		err = dir.Archive(targetTmpDir, filepath.Join(targetDir, m.ID+mtarSuffix))
		return err
	})
	if err != nil {
		return errors.Wrap(err, "MTAR Generation failed on MTA processing")
	}
	logs.Logger.Info("MTAR Generation successfully finished")
	return nil
}
