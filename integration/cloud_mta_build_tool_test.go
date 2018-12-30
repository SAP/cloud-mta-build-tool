package integration_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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
		cmd := exec.Command("go", "build", "-o", filepath.FromSlash("./integration/testdata/mtahtml5/"+mbtName), ".")
		cmd.Dir = filepath.FromSlash("../")
		err := cmd.Run()
		fmt.Println("finish to execute process", err)
		if err != nil {
			fmt.Println(err)
		}
	})

	AfterSuite(func() {
		os.Remove("./testdata/mtahtml5/" + mbtName)
		os.Remove("./testdata/mtahtml5/Makefile.mta")
		os.Remove("./testdata/mtahtml5/mtahtml5.mtar")
	})

	var _ = Describe("Command to provide the list of modules", func() {

		It("Getting module", func() {
			dir, _ := os.Getwd()
			args := "provide modules"

			path := dir + filepath.FromSlash("/testdata/mtahtml5")
			bin := filepath.FromSlash("./mbt")

			cmdOut, err := execute(bin, args, path)
			if len(err) > 0 {
				fmt.Println(err)
			}
			Ω(cmdOut).ShouldNot(BeNil())
			Ω(cmdOut).Should(BeEquivalentTo("[ui5app]" + "\n"))
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

			path := dir + filepath.FromSlash("/testdata/mtahtml5")
			bin := filepath.FromSlash("./mbt")

			cmdOut, err := execute(bin, args, path)
			if len(err) > 0 {
				fmt.Println(err)
			}
			Ω(cmdOut).ShouldNot(BeNil())

			//Read the MakeFile was generated
			out, error := ioutil.ReadFile(filepath.Join(dir, "testdata", "mtahtml5", "MakeFile.mta"))
			Ω(error).Should(Succeed())

			//Read the expected MakeFile
			expected, error := ioutil.ReadFile(filepath.Join(dir, "testdata", "ExpectedMakeFileWindows"))
			Ω(error).Should(Succeed())

			Ω(bytes.Equal(out, expected)).Should(BeTrue())
		})

		It("Command name error", func() {
			dir, _ := os.Getwd()
			args := "init 2"

			path := dir + filepath.FromSlash("/testdata/mtahtml5")
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
			fmt.Println(dir)
			path := dir + filepath.FromSlash("/testdata/mtahtml5")
			bin := filepath.FromSlash("make")
			cmdOut, err := execute(bin, args, path)
			if len(err) > 0 {
				fmt.Println(err)
			}
			Ω(err).Should(Equal(""))
			Ω(cmdOut).ShouldNot(BeEmpty())
		})
	})
})

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
