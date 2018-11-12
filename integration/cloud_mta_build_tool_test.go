package integration_test

import (
	"bytes"
	"fmt"
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

	BeforeEach(func() {
		By("Building MBT")
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			mbtName = "mbt"
		} else {
			mbtName = "mbt.exe"
		}
		cmd := exec.Command("go", "build", "-o", filepath.FromSlash("./integration/testdata/"+mbtName), ".")
		cmd.Dir = filepath.FromSlash("../")
		err := cmd.Run()
		fmt.Println("finish to execute process", err)
		if err != nil {
			fmt.Println(err)
		}
	})

	AfterEach(func() {
		os.Remove("./testdata/" + mbtName)
	})

	var _ = Describe("Command to provide the list of modules", func() {

		It("Getting module", func() {
			dir, _ := os.Getwd()
			args := "provide modules"

			path := dir + filepath.FromSlash("/testdata/")
			bin := filepath.FromSlash("./mbt")

			cmdOut, err := execute(bin, args, path)
			if len(err) > 0 {
				fmt.Println(err)
			}
			Ω(cmdOut).ShouldNot(BeNil())
			Ω(cmdOut).Should(BeEquivalentTo("[eb-java eb-db eb-ui-conf-eb eb-ui-conf-extensionfunction eb-ui-conf-movementcategory eb-ui-conf-stockledgercharacteristic eb-ui-conf-taxrate eb-ui-conf-taxwarehouse eb-ui-md-materialmaster eb-ui-md-shiptomaster eb-ui-stockledgerlineitem eb-ui-stockledgerlineitem-alp eb-ui-stockledgerprocessingerror eb-approuter eb-ftp-content eb-sb eb-msahaa eb-uideployer]" + "\n"))
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
})

// Execute commands and get outputs
func execute(bin string, args string, path string) (string, error string) {

	cmd := exec.Command(bin, strings.Split(args, " ")...)
	// bin path
	cmd.Dir = path
	// std out
	stdoutBuf := &bytes.Buffer{}
	cmd.Stdout = stdoutBuf
	// std error
	stdErrBuf := &bytes.Buffer{}
	cmd.Stderr = stdErrBuf

	cmd.Start()
	// wait to the command to finish
	err := cmd.Wait()
	if err != nil {
		fmt.Println(err)
	}
	return stdoutBuf.String(), stdErrBuf.String()
}
