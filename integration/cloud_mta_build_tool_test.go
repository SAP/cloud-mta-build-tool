// +build integration

package integration_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	archiveName = "mta_demo_0.0.1.mtar"
	binPath     = "mbt"
)

var _ = Describe("Integration - CloudMtaBuildTool", func() {

	var mbtName = ""

	BeforeSuite(func() {
		By("Building MBT")
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			mbtName = "mbt"
		} else {
			mbtName = "mbt.exe"
		}
		cmd := exec.Command("go", "build", "-o", filepath.Join(os.Getenv("GOPATH"), "/bin/"+mbtName), ".")
		cmd.Dir = filepath.FromSlash("../")
		err := cmd.Run()
		if err != nil {
			fmt.Println("binary creation failed: ", err)
		}
	})

	AfterSuite(func() {
		os.Remove("./testdata/mta_demo/" + mbtName)
		os.Remove("./testdata/mta_demo/Makefile.mta")
		os.Remove("./testdata/mta_demo/" + archiveName)
		resourceCleanup("node")
	})

	var _ = Describe("Command to provide the list of modules", func() {

		It("Getting module", func() {
			dir, _ := os.Getwd()
			path := dir + filepath.FromSlash("/testdata/mta_demo")
			bin := filepath.FromSlash(binPath)
			cmdOut, err, _ := execute(bin, "provide modules", path)
			Ω(err).Should(Equal(""))
			Ω(cmdOut).ShouldNot(BeNil())
			Ω(cmdOut).Should(ContainSubstring("[node]" + "\n"))
		})

		It("Command name error", func() {
			dir, _ := os.Getwd()
			path := dir + filepath.FromSlash("/testdata/")
			bin := filepath.FromSlash(binPath)
			_, err, _ := execute(bin, "provide modules 2", path)
			Ω(err).ShouldNot(BeNil())
		})
	})
	var _ = Describe("Generate the Makefile according to the mta.yaml file", func() {

		It("Generate Makefile", func() {
			dir, _ := os.Getwd()
			path := filepath.Join(dir, "testdata", "mta_demo")
			bin := filepath.FromSlash(binPath)
			cmdOut, err, _ := execute(bin, "init", path)
			if len(err) > 0 {
				fmt.Println(err)
			}
			Ω(cmdOut).ShouldNot(BeNil())
			// Read the MakeFile was generated
			out, error := ioutil.ReadFile(filepath.Join(dir, "testdata", "mta_demo", "Makefile.mta"))
			Ω(error).Should(BeNil())
			Ω(out).ShouldNot(BeEmpty())
		})

		It("Command name error", func() {
			dir, _ := os.Getwd()
			path := dir + filepath.FromSlash("/testdata/mta_demo")
			bin := filepath.FromSlash(binPath)
			_, err, _ := execute(bin, "init 2", path)
			Ω(err).ShouldNot(BeNil())
		})
	})

	var _ = Describe("Generate MTAR", func() {
		It("Generate MTAR", func() {

			dir, _ := os.Getwd()
			path := dir + filepath.FromSlash("/testdata/mta_demo")
			bin := filepath.FromSlash("make")
			cmdOut, err, _ := execute(bin, "-f Makefile.mta p=cf", path)
			if len(err) > 0 {
				fmt.Println(err)
			}
			Ω(err).Should(Equal(""))
			Ω(cmdOut).ShouldNot(BeEmpty())
			// Read the MakeFile was generated
			out, error := ioutil.ReadFile(filepath.Join(dir, "testdata", "mta_demo", "mta_archives", archiveName))
			Ω(error).Should(BeNil())
			Ω(out).ShouldNot(BeNil())
		})
	})

	var _ = Describe("Deploy basic mta archive", func() {
		It("Deploy MTAR", func() {
			dir, _ := os.Getwd()
			path := dir + filepath.FromSlash("/testdata/mta_demo/mta_archives")
			bin := filepath.FromSlash("cf")
			// Execute deployment process with output to make the deployment success/failure more clear
			executeWithOutput(bin, "deploy "+archiveName+" -f", path)
			// Check if the deploy succeeded by using curl command response.
			// Receiving the output status code 200 represents successful deployment
			args := "-s -o /dev/null -w '%{http_code}' " + os.Getenv("NODE_APP_ROUTE")
			path = dir + filepath.FromSlash("/testdata/mta_demo")
			bin = filepath.FromSlash("curl")
			cmdOut, err := executeEverySecond(bin, args, path)
			if len(err) > 0 {
				log.Println(err)
			}
			Ω(err).Should(Equal(""))
			Ω(cmdOut).Should(Equal("'200'"))
		})
	})
})

// execute with live output
func executeWithOutput(bin string, args string, path string) {
	cmd := exec.Command(bin, strings.Split(args, " ")...)
	cmd.Dir = path
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("process output | %s\n", scanner.Text())
		}
	}()
	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}
}

// Delete deployed app
func resourceCleanup(appName string) {
	dir, _ := os.Getwd()
	path := dir + filepath.FromSlash("/testdata/mta_demo")
	bin := filepath.FromSlash("cf")
	cmdOut, err, _ := execute(bin, "delete "+appName+" -r -f", path)
	if len(err) > 0 {
		fmt.Println(err)
	}
	Ω(err).Should(Equal(""))
	Ω(cmdOut).ShouldNot(BeEmpty())
}

// Execute command every second for 40 times
func executeEverySecond(bin string, args string, path string) (string, error string) {
	n := 0
	cmdOut, err, _ := execute(bin, args, path)
	for range time.Tick(time.Second) {
		cmdOut, err, _ = execute(bin, args, path)
		n++
		if n == 40 || strings.Compare(cmdOut, "'200'") == 0 {
			break
		}
	}
	return cmdOut, err
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
