package integration_test

import (
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

var _ = Describe("Integration - CloudMtaBuildTool", func() {

	var mbtName = ""

	BeforeSuite(func() {
		By("Building MBT")
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			mbtName = "mbt"
		} else {
			mbtName = "mbt.exe"
		}
		cmd := exec.Command("go", "build", "-o", filepath.FromSlash("./integration/testdata/mta_demo/"+mbtName), ".")
		cmd.Dir = filepath.FromSlash("../")
		err := cmd.Run()
		fmt.Println("finish to execute process", err)
		if err != nil {
			fmt.Println(err)
		}
	})

	AfterSuite(func() {
		os.Remove("./testdata/mta_demo/" + mbtName)
		os.Remove("./testdata/mta_demo/Makefile.mta")
		os.Remove("./testdata/mta_demo/mta_demo2.mtar")
		DeleteFromCF("node")
	})

	var _ = Describe("Command to provide the list of modules", func() {

		It("Getting module", func() {
			dir, _ := os.Getwd()
			args := "provide modules"

			path := dir + filepath.FromSlash("/testdata/mta_demo")
			bin := filepath.FromSlash("./mbt")

			cmdOut, err := execute(bin, args, path)
			if len(err) > 0 {
				fmt.Println(err)
			}
			Ω(cmdOut).ShouldNot(BeNil())
			Ω(cmdOut).Should(BeEquivalentTo("[node]" + "\n"))
		})

		It("Command name error", func() {
			dir, _ := os.Getwd()
			args := "provide modules 2"

			path := dir + filepath.FromSlash("/testdata/")
			bin := filepath.FromSlash("./mbt")

			cmdOut, err := execute(bin, args, path)
			if len(err) > 0 {
				fmt.Println(err)
			}
			Ω(err).ShouldNot(BeNil())
			Ω(cmdOut).Should(BeEmpty())
		})
	})
	var _ = Describe("Generate the Makefile according to the mta.yaml file", func() {

		It("Generate Makefile", func() {
			dir, _ := os.Getwd()
			args := "init"

			path := dir + filepath.FromSlash("/testdata/mta_demo")
			bin := filepath.FromSlash("./mbt")

			cmdOut, err := execute(bin, args, path)
			if len(err) > 0 {
				fmt.Println(err)
			}
			Ω(cmdOut).ShouldNot(BeNil())

			//Read the MakeFile was generated
			out, error := ioutil.ReadFile(filepath.Join(dir, "testdata", "mta_demo", "MakeFile.mta"))
			Ω(error).Should(Succeed())

			//Read the expected MakeFile
			expected, error := ioutil.ReadFile(filepath.Join(dir, "testdata", "ExpectedMakeFileWindows"))
			Ω(error).Should(Succeed())

			Ω(bytes.Equal(out, expected)).Should(BeTrue())
		})

		It("Command name error", func() {
			dir, _ := os.Getwd()
			args := "init 2"

			path := dir + filepath.FromSlash("/testdata/mta_demo")
			bin := filepath.FromSlash("./mbt")

			cmdOut, err := execute(bin, args, path)
			if len(err) > 0 {
				fmt.Println(err)
			}
			Ω(err).ShouldNot(BeNil())
			Ω(cmdOut).Should(BeEmpty())
		})
	})

	var _ = Describe("Generate MTAR", func() {
		It("Generate MTAR", func() {
			dir, _ := os.Getwd()
			args := "-f Makefile.mta p=cf"
			path := dir + filepath.FromSlash("/testdata/mta_demo")
			bin := filepath.FromSlash("make")
			cmdOut, err := execute(bin, args, path)

			if len(err) > 0 {
				fmt.Println(err)
			}
			Ω(err).Should(Equal(""))
			Ω(cmdOut).ShouldNot(BeEmpty())
		})
	})

	var _ = Describe("Deploy MTAR", func() {
		It("Deploy MTAR", func() {
			dir, _ := os.Getwd()
			args := "deploy mta_demo2.mtar"
			path := dir + filepath.FromSlash("/testdata/mta_demo")
			bin := filepath.FromSlash("cf")
			//Deploy Mtar
			cmdOut, err := polling(bin, args, path, 90)
			if len(err) > 0 {
				fmt.Println(err)
			}
			Ω(err).Should(Equal(""))
			Ω(cmdOut).ShouldNot(BeEmpty())

			//Command to check if the deploy succeeded by using curl command response.
			//After 90 seconds, there is a check if the deploy succeeded every second.
			//Receiving the output status code 200 represents the success.
			//If there is no success after 40 times, the test will fail.
			args = "-s -o /dev/null -w '%{http_code}' https://devx2-playg-node.cfapps.sap.hana.ondemand.com//"
			path = dir + filepath.FromSlash("/testdata/mta_demo")
			bin = filepath.FromSlash("curl")
			cmdOut, err = executeEverySecond(bin, args, path)
			if len(err) > 0 {
				fmt.Println(err)
			}
			Ω(err).Should(Equal(""))
			Ω(cmdOut).Should(Equal("'200'"))
		})
	})
})

// Delete app from cloud foundry
func DeleteFromCF (appName string) {
	dir, _ := os.Getwd()
	args := "delete " + appName + " -r -f"
	path := dir + filepath.FromSlash("/testdata/mta_demo")
	bin := filepath.FromSlash("cf")
	cmdOut, err := execute(bin, args, path)
	if len(err) > 0 {
		fmt.Println(err)
	}
	Ω(err).Should(Equal(""))
	Ω(cmdOut).ShouldNot(BeEmpty())
}

// Execute commands with wait timeout and get outputs
func polling(bin string, args string, path string, waitTimeOut time.Duration) (string, err string) {
	// Provide list of commands
	cmd := exec.Command(bin, strings.Split(args, " ")...)
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
	// Wait for the process to finish or kill it after a timeout (whichever happens first):
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(5 * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
		}
		log.Println("process killed as timeout reached")
	case err := <-done:
		if err != nil {
			log.Fatalf("process finished with error = %v", err)
		}
		log.Print("process finished successfully")
	}

	return stdoutBuf.String(), stdErrBuf.String()
}

// Execute command every second for 40 times
func executeEverySecond(bin string, args string, path string) (string, error string) {
	n := 0
	cmdOut, err := execute(bin, args, path)
		for range time.Tick(time.Second) {
			cmdOut, err = execute(bin, args, path)
		n++
		if n == 40 || strings.Compare(cmdOut, "'200'") == 0{
			break
		}
	}
	return cmdOut, err
}

// Execute commands and get outputs
func execute(bin string, args string, path string) (string, error string) {
	// Provide list of commands
	cmd := exec.Command(bin, strings.Split(args, " ")...)
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
	return stdoutBuf.String(), stdErrBuf.String()
}