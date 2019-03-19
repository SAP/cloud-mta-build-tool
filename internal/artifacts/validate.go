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
func ExecuteValidation(source, desc, mode, strict, exclude string, getWorkingDir func() (string, error)) error {
	logs.Logger.Info("validating the MTA project")

	strictValue, err := strconv.ParseBool(strict)
	if err != nil {
		return fmt.Errorf(`the "%s" strictness value is wrong; boolean value expected`, strict)
	}

	loc, err := dir.Location(source, "", desc, getWorkingDir)
	if err != nil {
		return errors.Wrap(err, "validation failed when initializing the location")
	}
	validateSchema, validateProject, err := validate.GetValidationMode(mode)
	if err != nil {
		return errors.Wrap(err, "validation failed when analyzing the validation mode")
	}
	warn, err := validate.MtaYaml(source, loc.GetMtaYamlFilename(), validateSchema, validateProject, strictValue, exclude)
	if warn != "" {
		logs.Logger.Warn(warn)
	}
	return err
}
