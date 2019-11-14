package exec

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kballard/go-shellquote"
	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/commands"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

func makeCommand(params []string) *exec.Cmd {
	if len(params) > 1 {
		// the command is running with user permission
		/* #nosec */
		return exec.Command(params[0], params[1:]...)
	}
	// the command is running with user permission
	/* #nosec */
	return exec.Command(params[0])
}

// ExecuteCommandsWithTimeout parses the list of commands and executes them in the current working directory with a specified timeout.
// If the timeout is reached an error is returned.
func ExecuteCommandsWithTimeout(commandsList []string, timeout string, path string, runIndicator bool) error {
	commandList, err := commands.CmdConverter(filepath.Clean(path), commandsList)
	if err != nil {
		return err
	}
	return ExecuteWithTimeout(commandList, timeout, runIndicator)
}

// ExecuteWithTimeout executes child processes and waits for the results. If the timeout is reached an error is returned and
// the child process is killed.
func ExecuteWithTimeout(cmdParams [][]string, timeout string, runIndicator bool) error {
	timeoutDuration, err := parseTimeoutString(timeout)
	if err != nil {
		return errors.Wrapf(err, ExecInvalidTimeoutMsg, timeout)
	}
	executeResultCh := make(chan error, 1)
	terminateCh := make(chan struct{})
	go func() {
		executeResultCh <- executeWithTerminateCh(cmdParams, terminateCh, runIndicator)
	}()

	select {
	case err = <-executeResultCh:
		return err
	case <-time.After(timeoutDuration):
		close(terminateCh)
		// Wait for executeWithTerminateCh to finish, to make sure we kill the running process
		err = <-executeResultCh
		if err != nil {
			logs.Logger.Error(err)
		}
		return errors.Errorf(ExecTimeoutMsg, timeoutDuration.String())
	}
}

func parseTimeoutString(timeoutString string) (time.Duration, error) {
	if timeoutString == "" {
		return 10 * time.Minute, nil
	}
	return time.ParseDuration(strings.TrimSpace(timeoutString))
}

// Execute - Execute child process and wait to results
func Execute(cmdParams [][]string, runIndicator bool) error {
	return executeWithTerminateCh(cmdParams, make(chan struct{}), runIndicator)
}

func executeWithTerminateCh(cmdParams [][]string, terminateCh <-chan struct{}, runIndicator bool) error {
	for _, cp := range cmdParams {
		var cmd *exec.Cmd
		commandString := shellquote.Join(cp[1:]...)
		logs.Logger.Infof(execMsg, commandString)
		cmd = makeCommand(cp[1:])
		cmd.Dir = cp[0]

		err := executeCommand(cmd, terminateCh, runIndicator)
		if err != nil {
			return errors.Wrapf(err, execFailed, commandString)
		}

	}
	return nil
}

// executeCommand - executes individual command
func executeCommand(cmd *exec.Cmd, terminateCh <-chan struct{}, runIndicator bool) error {
	logs.Logger.Debugf(execFileMsg, cmd.Path)

	// During the running process get the standard output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, execFailedOnStdoutMsg)
	}
	// During the running process get the standard output
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrap(err, execFailedOnStderrMsg)
	}

	// Start indicator if required
	if runIndicator {
		shutdownCh := make(chan struct{})
		go indicator(shutdownCh)
		defer close(shutdownCh) // Signal indicator() to terminate
	}

	// Start the process without waiting for it to finish
	if err = cmd.Start(); err != nil {
		return err
	}

	// Wait for the process to finish in a goroutine. We wait until it finishes or termination is requested via terminateCh.
	finishedCh := make(chan error, 1)
	go func() {
		// Stream command standard and error output.
		// Note: this does not wait until the process finishes, but will keep streaming its output until it is.
		pipeOutput(stdout)
		pipeOutput(stderr)

		// Get execution success or failure
		finishedCh <- cmd.Wait()
	}()

	select {
	case err = <-finishedCh:
		if err != nil {
			return err
		}
	case <-terminateCh:
		// Kill the process. We don't care if an error occurs here, we did our best and it doesn't affect the user.
		_ = cmd.Process.Kill()
		// Return an error so that we don't continue to the next process
		return fmt.Errorf(execKilledMsg)
	}

	return nil
}

// Scan the reader and write it to the output
func pipeOutput(reader io.Reader) {
	// instructs the scanner to read the input by runes instead of the default by-lines
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanRunes)
	go func() {
		for scanner.Scan() {
			fmt.Print(scanner.Text())
		}
	}()
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
