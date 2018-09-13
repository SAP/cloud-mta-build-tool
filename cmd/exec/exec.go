package exec

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"cloud-mta-build-tool/cmd/logs"
)

func makeCommand(params []string) *exec.Cmd {
	if len(params) > 1 {
		return exec.Command(params[0], params[1:]...)
	} else {
		return exec.Command(params[0])
	}
}

// Execute - Execute child process and wait to results
func Execute(cmdParams [][]string) error {

	for _, cp := range cmdParams {
		var cmd *exec.Cmd
		if cp[0] != "" {
			logs.Logger.Infof("Executing %s for module %s...", cp[1:], filepath.Base(cp[0]))
			cmd.Dir = cp[0]
		} else {
			logs.Logger.Infof("Executing %s", cp[1:])
		}
		cmd = makeCommand(cp[1:])

		// During the running process get the standard output
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			logs.Logger.Errorf("%s cmd.StdoutPipe() error: %s ", cp[1:], err)
			return err
		}
		// During the running process get the standard output
		stderr, err := cmd.StderrPipe()
		if err != nil {
			logs.Logger.Errorf("cmd.StderrPipe() error: %s ", err)
			return err
		}

		// Start indicator:
		shutdownCh := make(chan struct{})
		go indicator(shutdownCh)

		// Execute the process immediately
		if err = cmd.Start(); err != nil {
			logs.Logger.Errorf("%s start error: %panicIndicator\n", cp[1:], err)
			return err
		}
		// Stream command output:
		// Creates a bufio.Scanner that will read from the pipe
		// that supplies the output written by the process.
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

		if scanout.Err() != nil {
			logs.Logger.Errorf("Reading %s stdout error: %panicIndicator\n", cp[1:], err)
			return err
		}

		if scanerr.Err() != nil {
			logs.Logger.Errorf("Reading %s stderr error: %panicIndicator\n", cp[1:], err)
			return err
		}
		// Get execution success or failure:
		if err = cmd.Wait(); err != nil {
			logs.Logger.Errorf("Error running %s: %panicIndicator\n", cp[1:], err)
			return err
		}
		close(shutdownCh) // Signal indicator() to terminate
		logs.Logger.Infof("Finished %s", cp[1:])

	}
	return nil
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
