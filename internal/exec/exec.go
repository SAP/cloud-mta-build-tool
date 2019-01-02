package exec

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"cloud-mta-build-tool/internal/logs"
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
		if cp[0] != "" {
			logs.Logger.Infof("executing the %s command for the %s module ...", cp[1:], filepath.Base(cp[0]))
		} else {
			logs.Logger.Infof("executing the %s command ...", cp[1:])
		}
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
	logs.Logger.Infof("execution of the %v command started", cmd.Path)

	// During the running process get the standard output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrapf(err, "execution of the %v command failed when getting the stdout pipe", cmd.Path)
	}
	// During the running process get the standard output
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrapf(err, "execution of the %v command failed when getting the stderr pipe", cmd.Path)
	}

	// Start indicator
	shutdownCh := make(chan struct{})
	go indicator(shutdownCh)

	// Execute the process immediately
	if err = cmd.Start(); err != nil {
		return errors.Wrapf(err, "execution of the %v command failed when starting", cmd.Path)
	}
	// Stream command output:
	// Creates a bufio.Scanner that will read from the pipe
	// that supplies the output written by the process.
	scanout, scanerr := scanner(stdout, stderr)

	if scanerr.Err() != nil {
		return errors.Wrapf(err,
			"execution of the %v command failed when scanning the stdout and stderr pipes", cmd.Path)
	}

	if scanout.Err() != nil {
		return errors.Wrapf(err,
			"execution of the %v command failed when receiving an error from the scanout object", cmd.Path)
	}

	// Get execution success or failure:
	if err = cmd.Wait(); err != nil {
		return errors.Wrapf(err,
			"execution of the %v command failed while waiting for finish", cmd.Path)
	}
	close(shutdownCh) // Signal indicator() to terminate
	logs.Logger.Infof("execution of the %v command finished successfully", cmd.Path)
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
