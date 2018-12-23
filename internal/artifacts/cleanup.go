package artifacts

import (
	"os"

	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/internal/logs"
)

// ExecuteCleanup - cleanups temp artifacts
func ExecuteCleanup(source, target, desc string, wdGetter func() (string, error)) error {
	logs.Logger.Info("Cleanup started")
	// Remove temp folder
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrap(err, "Cleanup failed when initializing the location")
	}
	targetTmpDir := loc.GetTargetTmpDir()
	err = os.RemoveAll(targetTmpDir)
	if err != nil {
		return errors.Wrapf(err, "Cleanup failed when removing the <%v> folder", targetTmpDir)
	}
	logs.Logger.Info("Cleanup successfully finished")
	return nil
}
