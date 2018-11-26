package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"cloud-mta-build-tool/internal/fsys"
)

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

var _ = Describe("Provide", func() {
	It("Valid path to yaml", func() {

		out := executeAndProvideOutput(func() {
			Ω(provideModules(&dir.Loc{SourcePath: filepath.Join("testdata", "mtahtml5")})).Should(Succeed())
		})
		Ω(out).Should(ContainSubstring("[ui5app ui5app2]"))
	})

	It("Invalid path to yaml", func() {
		Ω(provideModules(&dir.Loc{SourcePath: filepath.Join("testdata", "mtahtml6")})).Should(HaveOccurred())
	})

	It("Invalid command call", func() {
		out := executeAndProvideOutput(func() {
			sourceBModuleFlag = ""
			pModuleCmd.RunE(nil, []string{})

		})
		Ω(out).Should(BeEmpty())
	})

})
