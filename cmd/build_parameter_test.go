package commands

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var source string

func copyFile(source string, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

var _ = Describe("mbt cli build to test build parameter", func() {
	BeforeEach(func() {
		source = "testdata/mtaignore"
		mbtCmdCLI = getBuildCmdCli()
	})

	AfterEach(func() {
		mbtCmdCLI = ""
		// Ω(os.RemoveAll(getTestPath("mtaignore", dir.MtarFolder))).Should(Succeed())
	})
	It("Success - build-parameter ignore all node_modules", func() {
		sourceMta := getTestPath("mtaignore", "mta-ignore-all-node_modules.yaml")
		targetMta := getTestPath("mtaignore", "mta.yaml")
		copyFile(sourceMta, targetMta)

		var stdout bytes.Buffer
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source)
		cmd.Stdout = &stdout

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.Remove(getTestPath("mtaignore", "mta.yaml"))).Should(Succeed())
	})
	It("Success - build-parameter ignore node_modules subfolders", func() {
		sourceMta := getTestPath("mtaignore", "mta-ignore-node_modules-subfolders.yaml")
		targetMta := getTestPath("mtaignore", "mta.yaml")
		copyFile(sourceMta, targetMta)

		var stdout bytes.Buffer
		cmd := exec.Command("bash", "-c", mbtCmdCLI+" build"+" --source "+source)
		cmd.Stdout = &stdout

		Ω(cmd.Run()).Should(Succeed())
		Ω(os.Remove(getTestPath("mtaignore", "mta.yaml"))).Should(Succeed())
	})
})
