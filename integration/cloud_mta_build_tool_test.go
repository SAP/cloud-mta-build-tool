// +build integration

package integration_test

import (
	"archive/zip"
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

	"github.com/SAP/cloud-mta/mta"
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
		os.Remove("./testdata/mta_demo/mtad.yaml")
		os.Remove("./testdata/mta_demo/" + archiveName)
		os.Remove("./testdata/mta_demo/mta_archives")
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

			// Check the MakeFile was generated
			Ω(filepath.Join(dir, "testdata", "mta_demo", "Makefile.mta")).Should(BeAnExistingFile())
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
			fmt.Println(cmdOut)
			Ω(cmdOut).ShouldNot(BeEmpty())
			// Check the archive was generated
			Ω(filepath.Join(dir, "testdata", "mta_demo", "mta_archives", archiveName)).Should(BeAnExistingFile())
		})
	})

	var _ = Describe("Generate the Verbose Makefile and use it for mtar generation", func() {

		It("Generate Verbose Makefile", func() {
			dir, _ := os.Getwd()
			os.RemoveAll(filepath.Join(dir, "testdata", "mta_demo", "Makefile.mta"))
			os.RemoveAll(filepath.Join(dir, "testdata", "mta_demo", "mta_archives", archiveName))
			path := filepath.Join(dir, "testdata", "mta_demo")
			bin := filepath.FromSlash(binPath)
			cmdOut, err, _ := execute(bin, "init -m=verbose", path)
			if len(err) > 0 {
				fmt.Println(err)
			}
			Ω(cmdOut).ShouldNot(BeNil())
			// Read the MakeFile was generated
			Ω(filepath.Join(dir, "testdata", "mta_demo", "Makefile.mta")).Should(BeAnExistingFile())
			// generate mtar
			bin = filepath.FromSlash("make")
			execute(bin, "-f Makefile.mta p=cf", path)
			// Check the archive was generated
			Ω(filepath.Join(dir, "testdata", "mta_demo", "mta_archives", archiveName)).Should(BeAnExistingFile())
		})

	})

	var _ = Describe("MBT gen commands", func() {
		It("Generate mtad", func() {
			dir, _ := os.Getwd()
			path := filepath.Join(dir, "testdata", "mta_demo")
			bin := filepath.FromSlash(binPath)
			_, err, _ := execute(bin, "gen mtad", path)
			Ω(err).Should(Equal(""))
			mtadPath := filepath.Join(path, "mtad.yaml")
			Ω(mtadPath).Should(BeAnExistingFile())
			content, _ := ioutil.ReadFile(mtadPath)
			mtadObj, _ := mta.Unmarshal(content)
			Ω(mtadObj.Modules[0].Type).Should(Equal("javascript.nodejs"))
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

	var _ = Describe("Assemble MTAR", func() {
		var currentWorkingDirectory string
		var mtaAssemblePath string
		It("Assemble MTAR", func() {
			currentWorkingDirectory, _ = os.Getwd()
			mtaAssemblePath = currentWorkingDirectory + filepath.FromSlash("/testdata/mta_assemble")

			bin := filepath.FromSlash(binPath)
			cmdOut, err, _ := execute(bin, "assemble", mtaAssemblePath)
			Ω(err).Should(Equal(""))
			Ω(cmdOut).ShouldNot(BeNil())
			Ω(cmdOut).Should(ContainSubstring("assembling the MTA project..." + "\n"))
			Ω(cmdOut).Should(ContainSubstring("copying the MTA content..." + "\n"))
			Ω(cmdOut).Should(ContainSubstring("generating the metadata..." + "\n"))
			Ω(cmdOut).Should(ContainSubstring("generating the MTA archive..." + "\n"))
			Ω(cmdOut).Should(ContainSubstring("the MTA archive generated at: " + filepath.Join(mtaAssemblePath, "mta_archives", "mta.assembly.example_1.3.3.mtar") + "\n"))
			Ω(cmdOut).Should(ContainSubstring("cleaning temporary files..." + "\n"))
			Ω(exists("mta.assembly.example_1.3.3.mtar", mtaAssemblePath+filepath.FromSlash("/mta_archives/"))).Should(BeTrue())
			validateMtaArchiveContents([]string{"META-INF/mtad.yaml", "META-INF/MANIFEST.MF", "node.zip", "xs-security.json"}, filepath.Join(mtaAssemblePath, "mta_archives", "mta.assembly.example_1.3.3.mtar"))
			os.Remove(filepath.Join(mtaAssemblePath, "mta.assembly.example.mtar"))
			os.Chdir(currentWorkingDirectory)
		})
	})
})

func validateMtaArchiveContents(expectedFilesInArchive []string, archiveLocation string) {
	archiveReader, err := zip.OpenReader(archiveLocation)
	Ω(err).Should(BeNil())
	defer archiveReader.Close()
	var filesInArchive []string
	for _, file := range archiveReader.File {
		filesInArchive = append(filesInArchive, file.Name)
	}
	for _, expectedFile := range expectedFilesInArchive {
		Ω(contains(expectedFile, filesInArchive)).Should(BeTrue())
	}

}

func contains(element string, elements []string) bool {
	for _, el := range elements {
		if el == element {
			return true
		}
	}
	return false
}

func exists(fileName, location string) bool {
	files, err := ioutil.ReadDir(location)
	Ω(err).Should(BeNil())
	for _, file := range files {
		if file.Name() == fileName {
			return true
		}
	}
	return false
}

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
