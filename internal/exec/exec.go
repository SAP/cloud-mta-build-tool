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
			logs.Logger.Infof("Executing %s for module %s...", cp[1:], filepath.Base(cp[0]))
		} else {
			logs.Logger.Infof("Executing %s", cp[1:])
		}
		cmd = makeCommand(cp[1:])
		cmd.Dir = cp[0]

		// During the running process get the standard output
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return errors.Wrapf(err, "%s cmd.StdoutPipe() error", cp[1:])
		}
		// During the running process get the standard output
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return errors.Wrapf(err, "%s cmd.StderrPipe() error", cp[1:])
		}

		// Start indicator
		shutdownCh := make(chan struct{})
		go indicator(shutdownCh)

		// Execute the process immediately
		if err = cmd.Start(); err != nil {
			return errors.Wrapf(err, "%s command start error", cp[1:])
		}
		// Stream command output:
		// Creates a bufio.Scanner that will read from the pipe
		// that supplies the output written by the process.
		scanout, scanerr := scanner(stdout, stderr)

		if scanout.Err() != nil {
			return errors.Wrapf(err, "%s scanout error", cp[1:])
		}

		if scanerr.Err() != nil {
			return errors.Wrapf(err, "Reading %s stderr error", cp[1:])
		}
		// Get execution success or failure:
		if err = cmd.Wait(); err != nil {
			return errors.Wrapf(err, "Error running %s", cp[1:])
		}
		close(shutdownCh) // Signal indicator() to terminate
		logs.Logger.Infof("Finished %s", cp[1:])

	}
	return nil
}

func scanner(stdout io.ReadCloser, stderr io.ReadCloser) (*bufio.Scanner, *bufio.Scanner) {
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
	// defer wg.Done()
	for {
		select {
		case <-ticker.C:
			fmt.Print(".")
		case <-shutdownCh:
			return
		}
	}
}
