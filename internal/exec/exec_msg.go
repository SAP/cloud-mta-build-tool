package exec

const (
	execMsg                 = `executing the "%s" command...`
	execFileMsg             = `the executable file is "%s"`
	execFailedOnStdoutMsg   = `could not execute the "%s" command when getting the stdout pipe`
	execFailedOnStderrMsg   = `could not execute the "%s" command when getting the stderr pipe`
	execFailedOnStartMsg    = `could not execute the "%s" command when starting`
	execFailedOnScanMsg     = `could not execute the "%s" command when scanning the stdout and stderr pipes`
	execFailedOnErrorGetMsg = `could not execute the "%s" command when receiving an error from the scanout object`
	execFailedOnFinishMsg   = `could not execute the "%s" command when waiting for finish`
	execInvalidTimeoutMsg   = `invalid timeout value "%s", it should be in the form "[123h][123m][123s]"`
	// ExecTimeoutMsg is the error message that occurs when a timeout is reached during commands execution
	ExecTimeoutMsg = `commands timed out after %s`
	execKilledMsg  = `process was terminated`
)
