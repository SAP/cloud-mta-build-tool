package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/tpl"
)

const (
	makefile    = "Makefile.mta"
	makefileTmp = "Makefile_tmp.mta"
)

// flags of init command
var initCmdSrc string
var initCmdTrg string
var initCmdDesc string
var initCmdName string
var initCmdMode string

// flags of build command
var buildProjectCmdSrc string
var buildProjectCmdTrg string
var buildProjectCmdDesc string
var buildProjectCmdName string
var buildProjectCmdMode string
var buildProjectCmdMtar string
var buildProjectCmdPlatform string

// init flags of init command
func init() {
	initCmd.Flags().StringVarP(&initCmdSrc, "source", "s", "", "Provide MTA source")
	initCmd.Flags().StringVarP(&initCmdTrg, "target", "t", "", "Provide MTA target")
	initCmd.Flags().StringVarP(&initCmdDesc, "desc", "d", "", "Descriptor MTA - dev/dep")
	initCmd.Flags().StringVarP(&initCmdDesc, "name", "n", "", "Name of Makefile")
	initCmd.Flags().StringVarP(&initCmdMode, "mode", "m", "", "Mode of Makefile generation - default/verbose")

	buildCmd.Flags().StringVarP(&buildProjectCmdSrc, "source", "s", "$(CURDIR)", "Provide MTA source")
	buildCmd.Flags().StringVarP(&buildProjectCmdTrg, "target", "t", "$(CURDIR)/mta_archives", "Provide MTA target")
	buildCmd.Flags().StringVarP(&buildProjectCmdDesc, "desc", "d", "", "Descriptor MTA - dev/dep")
	buildCmd.Flags().StringVarP(&buildProjectCmdName, "name", "n", "", "Name of Makefile")
	buildCmd.Flags().StringVarP(&buildProjectCmdMode, "mode", "m", "", "Name of Mtar")
	buildCmd.Flags().StringVarP(&buildProjectCmdMtar, "mtar", "", "*", "Mode of Makefile generation - default/verbose")
	buildCmd.Flags().StringVarP(&buildProjectCmdPlatform, "platform", "p", "", "Platform")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generates Makefile",
	Long:  "generates Makefile as a manifest file that describes the build process",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Generate build script
		err := tpl.ExecuteMake(initCmdSrc, initCmdTrg, makefile, initCmdDesc, initCmdMode, os.Getwd)
		logError(err)
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "generates and executes Makefile",
	Long:  "generates Makefile as a manifest file that describes the build process and executes it",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Generate build script
		err := tpl.ExecuteMake(buildProjectCmdSrc, "", makefileTmp, buildProjectCmdDesc, buildProjectCmdMode, os.Getwd)
		logError(err)
		bin := filepath.FromSlash("make")
		commandArgs := "-f " + makefileTmp + " p=" + buildProjectCmdPlatform + " mtar=" + buildProjectCmdMtar + ` t="` + buildProjectCmdTrg + `"`
		cmdout, error, _ := execute(bin, commandArgs, buildProjectCmdSrc)
		fmt.Println(cmdout)
		if error != "" {
			fmt.Println("binary creation failed: ", err)
		}
		logError(err)
		// Remove Makefile_tmp.mta file from directory
		err = os.Remove(filepath.FromSlash(makefileTmp))
		if err != nil {
			fmt.Println("deleting"+makefileTmp+"", err)
		}
	},
}

// Execute commands and get outputs
func execute(bin string, args string, path string) (string, error string, cmd *exec.Cmd) {
	// Provide list of commands
	cmd = exec.Command(bin, strings.Split(args, " ")...)
	// bin path
	cmd.Dir = path
	// std out
	stdoutBuf := &bytes.Buffer{}
	cmd.Stdout = stdoutBuf
	// std error
	stdErrBuf := &bytes.Buffer{}
	cmd.Stderr = stdErrBuf
	// Start command
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}
	// wait to the command to finish
	err := cmd.Wait()
	if err != nil {
		fmt.Println(err)
	}
	return stdoutBuf.String(), stdErrBuf.String(), cmd
}
