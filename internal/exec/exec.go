package exec

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
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

// Execute - Execute child process and wait to results
func Execute(cmdParams [][]string) error {

	for _, cp := range cmdParams {
		var cmd *exec.Cmd
		logs.Logger.Infof(execMsg, cp[1:])
		cmd = makeCommand(cp[1:])
		cmd.Dir = cp[0]

		err := executeCommand(cmd)
		if err != nil {
			return err
		}

	}
	return nil
}

// executeCommand - executes individual command
func executeCommand(cmd *exec.Cmd) error {
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

	// Execute the process immediately
	if err = cmd.Start(); err != nil {
		return errors.Wrapf(err, execFailedOnStartMsg, cmd.Path)
	}
	// Stream command output:
	// Creates a bufio.Scanner that will read from the pipe
	// that supplies the output written by the process.
	scanout, scanerr := scanner(stdout, stderr)

	if scanerr.Err() != nil {
		return errors.Wrapf(err, execFailedOnScanMsg, cmd.Path)
	}

	if scanout.Err() != nil {
		return errors.Wrapf(err, execFailedOnErrorGetMsg, cmd.Path)
	}

	// Get execution success or failure:
	if err = cmd.Wait(); err != nil {
		return errors.Wrapf(err, execFailedOnFinishMsg, cmd.Path)
	}
	close(shutdownCh) // Signal indicator() to terminate
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
