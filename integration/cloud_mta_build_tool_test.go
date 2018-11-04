package integration_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudMtaBuildTool", func() {

	BeforeEach(func() {
		By("Building MBT")
		cmd := exec.Command("go", "build", "-o", "./integration/testdata/mbt", ".")
		cmd.Dir = "../"
		err := cmd.Run()
		fmt.Println("finish to execute process", err)
		if err != nil {
			fmt.Println(err)
		}
	})

	AfterEach(func() {
		os.Remove("./testdata/mbt")
	})

	var _ = Describe("Command to provide the list of modules", func() {
		It("Target file in opened status", func() {
			dir,_ := os.Getwd()
			path := dir + "/testdata/"
			args := "provide modules"
			bin := "./mbt"
			cmdOut := execute(bin, args, path)
			Ω(cmdOut).ShouldNot(BeNil())
			Ω(cmdOut).Should(BeEquivalentTo("[eb-java eb-db eb-uideployer eb-ui-conf-eb eb-ui-conf-extensionfunction eb-ui-conf-movementcategory eb-ui-conf-stockledgercharacteristic eb-ui-conf-taxrate eb-ui-conf-taxwarehouse eb-ui-md-materialmaster eb-ui-md-shiptomaster eb-ui-stockledgerlineitem eb-ui-stockledgerlineitem-alp eb-ui-stockledgerprocessingerror eb-approuter eb-ftp-content eb-sb eb-msahaa]" + "\n"))
		})

	})
})

// Execute commands and get outputs
func execute(bin string, args string, path string) string {

	cmd := exec.Command(bin, strings.Split(args, " ")...)
	// bin path
	cmd.Dir = path
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanRunes)
	var buffer bytes.Buffer
	for scanner.Scan() {
		fmt.Printf(scanner.Text())
		buffer.WriteString(scanner.Text())
	}
	// wait to the command to finish
	err := cmd.Wait()
	if err != nil {
		fmt.Println(err)
	}
	return buffer.String()
}
