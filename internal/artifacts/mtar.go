package artifacts

import (
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"strconv"
)

const (
	mtarExtension = ".mtar"
)

// ExecuteGenMtar - generates MTAR
func ExecuteGenMtar(source, target, targetProvided, desc string, wdGetter func() (string, error)) error {
	logs.Logger.Info("generating the MTA archive...")
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrap(err, "generation of the MTA archive failed when initializing the location")
	}
	path, err := generateMtar(loc, loc, loc, isTargetProvided(source, target, targetProvided))
	if err != nil {
		return err
	}
	logs.Logger.Infof("the MTA archive generated at: %s", path)
	return nil
}

func isTargetProvided(source, target, provided string) bool {
	if provided == "" {
		return target != ""
	}
	value, err := strconv.ParseBool(provided)
	if err != nil {
		return false
	}
	return value
}

// generateMtar - generate mtar archive from the build artifacts
func generateMtar(targetLoc dir.ITargetPath, targetArtifacts dir.ITargetArtifacts, parser dir.IMtaParser, targetProvided bool) (string, error) {
	// get MTA object
	m, err := parser.ParseFile()
	if err != nil {
		return "", errors.Wrap(err, "generation of the the MTA archive failed when parsing the mta file")
	}
	// get target temporary folder to be archived
	targetTmpDir := targetLoc.GetTargetTmpDir()

	// create the mta_archives folder
	// get directory - where mtar will be saved
	mtarFolderPath := targetArtifacts.GetMtarDir(targetProvided)
	err = dir.CreateDirIfNotExist(mtarFolderPath)
	if err != nil {
		return "", errors.Wrapf(err,
			`generation of the MTA archive failed when creating the "%s" folder`, mtarFolderPath)
	}
	// archive building artifacts to mtar
	mtarPath := filepath.Join(mtarFolderPath, m.ID+"_"+m.Version+mtarExtension)
	err = dir.Archive(targetTmpDir, mtarPath)
	if err != nil {
		return "", errors.Wrap(err, "generation of the MTA archive failed when archiving")
	}
	return mtarPath, nil
}
