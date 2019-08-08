package artifacts

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

// Assembly - assemble mta project
func Assembly(source, target string, extensions []string, platform, mtarName, copyInParallel string, getWd func() (string, error)) error {

	logs.Logger.Info(assemblingMsg)

	parallelCopy, err := strconv.ParseBool(copyInParallel)
	if err != nil {
		parallelCopy = false
	}
	// copy from source to target
	err = CopyMtaContent(source, target, extensions, parallelCopy, getWd)
	if err != nil {
		return errors.Wrap(err, assemblyFailedOnCopyMsg)
	}
	// Generate meta artifacts
	err = ExecuteGenMeta(source, target, dir.Dep, extensions, platform, getWd)
	if err != nil {
		return errors.Wrap(err, assemblyFailedOnMetaMsg)
	}
	// generate mtar
	err = ExecuteGenMtar(source, target, strconv.FormatBool(target != ""), dir.Dep, extensions, mtarName, getWd)
	if err != nil {
		return errors.Wrap(err, assemblyFailedOnMtarMsg)
	}
	// cleanup
	err = ExecuteCleanup(source, target, dir.Dep, getWd)
	if err != nil {
		return errors.Wrap(err, assemblyFailedOnCleanupMsg)
	}
	return nil
}
