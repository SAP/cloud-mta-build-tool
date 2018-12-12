package artifacts

import (
	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/internal/logs"
	"cloud-mta-build-tool/validations"
)

// ExecuteValidation - executes validation of MTA
func ExecuteValidation(source, desc, mode string, getWorkingDir func() (string, error)) error {
	logs.Logger.Info("MBT Validation started")
	loc, err := dir.Location(source, "", desc, getWorkingDir)
	if err != nil {
		return errors.Wrap(err, "MBT Validation failed on location initialization")
	}
	validateSchema, validateProject, err := validate.GetValidationMode(mode)
	if err != nil {
		return errors.Wrap(err, "MBT Validation failed on validation mode analysis")
	}
	err = validate.ValidateMtaYaml(source, loc.GetMtaYamlFilename(), validateSchema, validateProject)
	if err != nil {
		return errors.Wrap(err, "MBT Validation failed")
	}
	logs.Logger.Info("MBT Validation successfully finished")
	return nil
}
