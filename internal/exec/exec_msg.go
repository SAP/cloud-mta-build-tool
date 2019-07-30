package exec

const (
	execMsg                = `executing the "%s" command...`
	execFileMsg            = `the executable file is at "%s"`
	execFailed             = `could not execute the "%s" command`
	execFailedOnStdoutMsg  = `could not get the "stdout" pipe`
	execFailedOnStderrMsg  = `could not get the "stderr" pipe`
	execFailedOnScanerrMsg = `could not read from the "stderr" pipe`
	execFailedOnScanoutMsg = `could not read from the "stdout" pipe`
	// ExecInvalidTimeoutMsg is the error message that occurs when a timeout value is invalid
	ExecInvalidTimeoutMsg = `invalid timeout value "%s", it should be in the form "[123h][123m][123s]"`
	// ExecTimeoutMsg is the error message that occurs when a timeout is reached during commands execution
	ExecTimeoutMsg = `the build timed out after %s`
	execKilledMsg  = `the process was interrupted`
)
