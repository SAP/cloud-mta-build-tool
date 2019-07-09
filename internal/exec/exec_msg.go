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
)
