package artifacts

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

// Assembly - assemble mta project
func Assembly(source, target, platform, mtarName, copyInParallel string, getWd func() (string, error)) error {

	logs.Logger.Info("assembling the MTA project...")

	parallelCopy, err := strconv.ParseBool(copyInParallel)
	if err != nil {
		parallelCopy = false
	}
	// copy from source to target
	err = CopyMtaContent(source, target, parallelCopy, getWd)
	if err != nil {
		return errors.Wrap(err, "assembly of the MTA project failed when copying the MTA content")
	}
	// Generate meta artifacts
	err = ExecuteGenMeta(source, target, dir.Dep, platform, getWd)
	if err != nil {
		return errors.Wrap(err, "assembly of the MTA project failed when generating the meta information")
	}
	// generate mtar
	err = ExecuteGenMtar(source, target, strconv.FormatBool(target != ""), dir.Dep, mtarName, getWd)
	if err != nil {
		return errors.Wrap(err, "assembly of the MTA project failed when generating the MTA archive")
	}
	// cleanup
	err = ExecuteCleanup(source, target, dir.Dep, getWd)
	if err != nil {
		return errors.Wrap(err, "assembly of the MTA project failed when executing cleanup")
	}
	return nil
}
