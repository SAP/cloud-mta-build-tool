package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var _ = BeforeSuite(func() {
	logs.Logger = logs.NewLogger()
})

func executeAndProvideOutput(execute func()) string {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	execute()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		if err != nil {
			fmt.Println(err)
		}
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC
	return out
}

func createFileInTmpFolder(projectName string, path ...string) {
	file, err := os.Create(getFullPathInTmpFolder(projectName, path...))
	Ω(err).Should(Succeed())
	err = file.Close()
	Ω(err).Should(Succeed())
}

func createDirInTmpFolder(projectName string, path ...string) {
	err := dir.CreateDirIfNotExist(getFullPathInTmpFolder(projectName, path...))
	Ω(err).Should(Succeed())
}

func getFullPathInTmpFolder(projectName string, path ...string) string {
	pathWithResultFolder := []string{"result", "." + projectName + "_mta_build_tmp"}
	pathWithResultFolder = append(pathWithResultFolder, path...)
	return getTestPath(pathWithResultFolder...)
}
