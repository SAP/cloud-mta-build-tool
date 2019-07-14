package exec

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

func makeCommand(params []string) *exec.Cmd {
	if len(params) > 1 {
		return exec.Command(params[0], params[1:]...)
	}
	return exec.Command(params[0])
}

// ExecuteWithTimeout executes child processes and waits for the results. If the timeout is reached an error is returned and
// the child process is killed.
func ExecuteWithTimeout(cmdParams [][]string, timeout string) error {
	timeoutDuration, err := parseTimeoutString(timeout)
	if err != nil {
		return errors.Wrapf(err, execInvalidTimeoutMsg, timeout)
	}
	executeResultCh := make(chan error, 1)
	terminateCh := make(chan struct{})
	go func() {
		executeResultCh <- executeWithTerminateCh(cmdParams, terminateCh)
	}()

	select {
	case err = <-executeResultCh:
		return err
	case <-time.After(timeoutDuration):
		close(terminateCh)
		// Wait for executeWithTerminateCh to finish, to make sure we kill the running process
		<-executeResultCh
		return errors.Errorf(ExecTimeoutMsg, timeoutDuration.String())
	}
}

func parseTimeoutString(timeoutString string) (time.Duration, error) {
	if timeoutString == "" {
		return 5 * time.Minute, nil
	}
	return time.ParseDuration(strings.TrimSpace(timeoutString))
}

// Execute - Execute child process and wait to results
func Execute(cmdParams [][]string) error {
	return executeWithTerminateCh(cmdParams, make(chan struct{}))
}

func executeWithTerminateCh(cmdParams [][]string, terminateCh <-chan struct{}) error {
	for _, cp := range cmdParams {
		var cmd *exec.Cmd
		logs.Logger.Infof(execMsg, cp[1:])
		cmd = makeCommand(cp[1:])
		cmd.Dir = cp[0]

		err := executeCommand(cmd, terminateCh)
		if err != nil {
			return err
		}

	}
	return nil
}

// executeCommand - executes individual command
func executeCommand(cmd *exec.Cmd, terminateCh <-chan struct{}) error {
	logs.Logger.Debugf(execFileMsg, cmd.Path)

	// During the running process get the standard output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrapf(err, execFailedOnStdoutMsg, cmd.Path)
	}
	// During the running process get the standard output
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrapf(err, execFailedOnStderrMsg, cmd.Path)
	}

	// Start indicator
	shutdownCh := make(chan struct{})
	go indicator(shutdownCh)
	defer close(shutdownCh) // Signal indicator() to terminate

	// Start the process without waiting for it to finish
	if err = cmd.Start(); err != nil {
		return errors.Wrapf(err, execFailedOnStartMsg, cmd.Path)
	}

	// Wait for the process to finish in a goroutine. We wait until it finishes or termination is requested via terminateCh.
	finishedCh := make(chan error, 1)
	go func() {
		// Stream command output:
		// Creates a bufio.Scanner that will read from the pipe
		// that supplies the output written by the process.
		// Note: this waits until the process finishes or an error occurs.
		scanout, scanerr := scanner(stdout, stderr)

		if scanerr.Err() != nil {
			finishedCh <- errors.Wrapf(err, execFailedOnScanMsg, cmd.Path)
			return
		}

		if scanout.Err() != nil {
			finishedCh <- errors.Wrapf(err, execFailedOnErrorGetMsg, cmd.Path)
			return
		}

		// Get execution success or failure
		finishedCh <- cmd.Wait()
	}()

	select {
	case err = <-finishedCh:
		if err != nil {
			return errors.Wrapf(err, execFailedOnFinishMsg, cmd.Path)
		}
	case <-terminateCh:
		// Kill the process. We don't care if an error occurs here, we did our best and it doesn't affect the user.
		_ = cmd.Process.Kill()
		// Return an error so that we don't continue to the next process
		return fmt.Errorf(execKilledMsg)
	}

	return nil
}

func scanner(stdout io.Reader, stderr io.Reader) (*bufio.Scanner, *bufio.Scanner) {
	scanout := bufio.NewScanner(stdout)
	scanerr := bufio.NewScanner(stderr)
	// instructs the scanner to read the input by runes instead of the default by-lines.
	scanout.Split(bufio.ScanRunes)
	for scanout.Scan() {
		fmt.Print(scanout.Text())
	}
	// instructs the scanner to read the input by runes instead of the default by-lines.
	scanerr.Split(bufio.ScanRunes)
	for scanerr.Scan() {
		fmt.Print(scanerr.Text())
	}
	return scanout, scanerr
}

// Show progress when the command is executed
// and the terminal are not providing any process feedback
func indicator(shutdownCh <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Print(".")
		case <-shutdownCh:
			return
		}
	}
}
