package artifacts

import (
	"path/filepath"

	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/internal/logs"
)

const (
	mtarExtension = ".mtar"
	mtarFolder    = "mta_archives"
)

// ExecuteGenMtar - generates MTAR
func ExecuteGenMtar(source, target, desc string, wdGetter func() (string, error)) error {
	logs.Logger.Info("generating the MTA archive")
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrap(err, "generation of the MTA archive failed when initializing the location")
	}
	path, err := generateMtar(loc, loc)
	if err != nil {
		return err
	}
	logs.Logger.Info("the MTA archive generated at: ", path)
	return nil
}

// generateMtar - generate mtar archive from the build artifacts
func generateMtar(targetLoc dir.ITargetPath, parser dir.IMtaParser) (string, error) {
	// get MTA object
	m, err := parser.ParseFile()
	if err != nil {
		return "", errors.Wrap(err, "generation of the the MTA archive failed when parsing the mta file")
	}
	// get target temporary folder to be archived
	targetTmpDir := targetLoc.GetTargetTmpDir()

	// create the mta_archives folder
	// get directory - where mtar will be saved
	mtarFolderPath := filepath.Join(targetLoc.GetTarget(), mtarFolder)
	err = dir.CreateDirIfNotExist(mtarFolderPath)
	if err != nil {
		return "", errors.Wrap(err, "generation of the MTA archive failed when creating the mta_archives folder")
	}
	// archive building artifacts to mtar
	mtarPath := filepath.Join(mtarFolderPath, m.ID+"_"+m.Version+mtarExtension)
	err = dir.Archive(targetTmpDir, mtarPath)
	if err != nil {
		return "", errors.Wrap(err, "generation of the MTA archive failed when archiving")
	}
	return mtarPath, nil
}
