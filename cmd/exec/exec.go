package exec

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"mbtv2/cmd/logs"
)

func Execute(cmdParams [][]string) error {

	for _, cp := range cmdParams {

		logs.Logger.Infof("Executing %s for module %s...", cp[1:], filepath.Base(cp[0]))
		cmd := exec.Command(cp[1], cp[2:]...)

		cmd.Dir = cp[0]
		// During the running process get the standard output
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			logs.Logger.Errorf("%s cmd.StdoutPipe() error: %s ", cp[1:], err, "\n")
			return err
		}
		// During the running process get the standard output
		stderr, err := cmd.StderrPipe()
		if err != nil {
			logs.Logger.Errorf("cmd.StderrPipe() error: ", cp[1:], err, "\n")
			return err
		}

		// Start indicator:
		shutdownCh := make(chan struct{})
		go indicator(shutdownCh)

		// Execute the process immediately
		if err = cmd.Start(); err != nil {
			logs.Logger.Errorf("%s start error: %v\n", cp[1:], err)
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
			logs.Logger.Errorf("Reading %s stdout error: %v\n", cp[1:], err)
			return err
		}

		if scanerr.Err() != nil {
			logs.Logger.Errorf("Reading %s stderr error: %v\n", cp[1:], err)
			return err
		}
		// Get execution success or failure:
		if err = cmd.Wait(); err != nil {
			logs.Logger.Errorf("Error running %s: %v\n", cp[1:], err)
			return err
		}
		close(shutdownCh) // Signal indicator() to terminate
		logs.Logger.Infof("Finished %s", cp[1:])

	}
	return nil
}

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
