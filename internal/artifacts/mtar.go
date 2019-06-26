package artifacts

import (
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"github.com/SAP/cloud-mta/mta"
	"strconv"
)

const (
	mtarExtension = ".mtar"
)

// ExecuteGenMtar - generates MTAR
func ExecuteGenMtar(source, target, targetProvided, desc, mtarName string, wdGetter func() (string, error)) error {
	logs.Logger.Info("generating the MTA archive...")
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrap(err, "generation of the MTA archive failed when initializing the location")
	}
	path, err := generateMtar(loc, loc, loc, isTargetProvided(target, targetProvided), mtarName)
	if err != nil {
		return err
	}
	logs.Logger.Infof("the MTA archive generated at: %s", path)
	return nil
}

func isTargetProvided(target, provided string) bool {
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
func generateMtar(targetLoc dir.ITargetPath, targetArtifacts dir.ITargetArtifacts, parser dir.IMtaParser,
	targetProvided bool, mtarName string) (string, error) {
	// get MTA object
	m, err := parser.ParseFile()
	if err != nil {
		return "", errors.Wrap(err, "generation of the the MTA archive failed when parsing the mta file")
	}
	// get target temporary folder to be archived
	targetTmpDir := targetLoc.GetTargetTmpDir()

	// get directory - where mtar will be saved
	mtarFolderPath := targetArtifacts.GetMtarDir(targetProvided)

	// archive building artifacts to mtar
	mtarPath := filepath.Join(mtarFolderPath, getMtarFileName(m, mtarName))
	err = dir.Archive(targetTmpDir, mtarPath, nil)
	if err != nil {
		return "", errors.Wrap(err, "generation of the MTA archive failed when archiving")
	}
	return mtarPath, nil
}

func getMtarFileName(m *mta.MTA, mtarName string) string {
	if mtarName == "" || mtarName == "*" {
		return m.ID + "_" + m.Version + mtarExtension
	}
	if filepath.Ext(mtarName) != "" {
		return mtarName
	}
	return mtarName + mtarExtension
}
