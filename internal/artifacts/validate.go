package artifacts

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	validate "github.com/SAP/cloud-mta/validations"
)

// ExecuteValidation - executes validation of MTA
func ExecuteValidation(source, mtaYamlFilename, desc string, extensions []string, mode, strict, exclude string, getWorkingDir func() (string, error)) error {
	logs.Logger.Info(validationMsg)

	strictValue, err := strconv.ParseBool(strict)
	if err != nil {
		return fmt.Errorf(wrongStrictnessMsg, strict)
	}

	loc, err := dir.Location(source, mtaYamlFilename, "", desc, extensions, getWorkingDir)
	if err != nil {
		return errors.Wrap(err, validationFailedOnLocMsg)
	}
	validateSchema, validateSemantic, err := validate.GetValidationMode(mode)
	if err != nil {
		return errors.Wrap(err, validationFailedOnModeMsg)
	}

	var validationErrors []string
	warn, err := validate.MtaYaml(source, loc.GetMtaYamlFilename(), validateSchema, validateSemantic, strictValue, exclude)
	if warn != "" {
		logs.Logger.Warnf("%s: %s", loc.GetMtaYamlFilename(), warn)
	}
	if err != nil {
		validationErrors = append(validationErrors, err.Error())
	}

	for _, ext := range loc.GetExtensionFilePaths() {
		warn, err = validate.Mtaext(source, ext, validateSchema, validateSemantic, strictValue, exclude)
		if warn != "" {
			logs.Logger.Warnf("%s: %s", ext, warn)
		}
		if err != nil {
			validationErrors = append(validationErrors, err.Error())
		}
	}

	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, "\n"))
	}

	return nil
}
