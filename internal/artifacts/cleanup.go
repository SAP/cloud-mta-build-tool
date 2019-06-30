package artifacts

import (
	"os"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

// ExecuteCleanup - cleanups temp artifacts
func ExecuteCleanup(source, target, desc string, wdGetter func() (string, error)) error {
	logs.Logger.Info("cleaning temporary files...")
	// Remove temp folder
	loc, err := dir.Location(source, target, desc, wdGetter)
	if err != nil {
		return errors.Wrap(err, "cleanup failed when initializing the location")
	}
	targetTmpDir := loc.GetTargetTmpDir()
	err = os.RemoveAll(targetTmpDir)
	if err != nil {
		return errors.Wrapf(err, `cleanup failed when removing the "%s" folder`, targetTmpDir)
	}
	return nil
}
