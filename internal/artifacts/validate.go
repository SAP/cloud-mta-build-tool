package artifacts

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"github.com/SAP/cloud-mta/validations"
)

// ExecuteValidation - executes validation of MTA
func ExecuteValidation(source, desc string, extensions []string, mode, strict, exclude string, getWorkingDir func() (string, error)) error {
	logs.Logger.Info(validationMsg)

	strictValue, err := strconv.ParseBool(strict)
	if err != nil {
		return fmt.Errorf(wrongStrictnessMsg, strict)
	}

	loc, err := dir.Location(source, "", desc, extensions, getWorkingDir)
	if err != nil {
		return errors.Wrap(err, validationFailedOnLocMsg)
	}
	validateSchema, validateProject, err := validate.GetValidationMode(mode)
	if err != nil {
		return errors.Wrap(err, validationFailedOnModeMsg)
	}
	warn, err := validate.MtaYaml(source, loc.GetMtaYamlFilename(), validateSchema, validateProject, strictValue, exclude)
	if warn != "" {
		logs.Logger.Warn(warn)
	}
	return err
}
