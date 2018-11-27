package artifacts

import (
	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fsys"
	"cloud-mta-build-tool/internal/logs"
)

// process mta.yaml file
func processMta(processName string, ep *dir.Loc, args []string, process func(file []byte, args []string) error) error {
	logs.Logger.Info("Starting " + processName)
	mf, err := dir.Read(ep)
	if err == nil {
		err = process(mf, args)
		if err == nil {
			logs.Logger.Info(processName + " finish successfully ")
		}
	} else {
		err = errors.Wrap(err, "MTA file not found")
	}
	return err
}
