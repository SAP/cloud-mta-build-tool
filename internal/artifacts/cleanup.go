package artifacts

import (
	"os"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

// ExecuteCleanup - cleanups temp artifacts
func ExecuteCleanup(source, target, desc string, wdGetter func() (string, error)) error {
	logs.Logger.Info(cleanupMsg)
	// Remove temp folder
	loc, err := dir.Location(source, target, desc, nil, wdGetter)
	if err != nil {
		return errors.Wrap(err, cleanupFailedOnLocMsg)
	}
	targetTmpDir := loc.GetTargetTmpDir()
	err = os.RemoveAll(targetTmpDir)
	if err != nil {
		return errors.Wrapf(err, cleanupFailedOnFolderMsg, targetTmpDir)
	}
	return nil
}
