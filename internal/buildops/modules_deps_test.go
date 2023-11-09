package buildops

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"strings"

	dir "github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta/mta"
)

var _ = Describe("ModulesDeps", func() {

	var _ = Describe("Process Dependencies test", func() {
		AfterEach(func() {
			os.RemoveAll(getTestPath("result"))
		})

		It("Sanity", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result"), MtaFilename: "mtaWithBuildParams.yaml"}
			Ω(ProcessDependencies(&ep, &ep, "ui5app")).Should(Succeed())
		})
		It("Invalid artifacts", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getTestPath("result"), MtaFilename: "mtaWithBuildParamsWithWrongArtifacts.yaml"}
			Ω(ProcessDependencies(&ep, &ep, "ui5app")).Should(HaveOccurred())
		})
		It("Invalid mta", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), MtaFilename: "mta1.yaml"}
			Ω(ProcessDependencies(&ep, &ep, "ui5app")).Should(HaveOccurred())
		})
		It("Invalid module name", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5")}
			Ω(ProcessDependencies(&ep, &ep, "xxx")).Should(HaveOccurred())
		})
		It("Invalid module name", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), MtaFilename: "mtaWithWrongBuildParams.yaml"}
			Ω(ProcessDependencies(&ep, &ep, "ui5app")).Should(HaveOccurred())
		})
	})

	It("Resolve dependencies - Valid case", func() {
		wd, _ := os.Getwd()
		ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mta_multiapps.yaml"}
		mtaStr, _ := ep.ParseFile(true)
		actual, _ := getModulesOrder(mtaStr)
		// last module depends on others
		Ω(actual[len(actual)-1]).Should(Equal("eb-uideployer"))
	})

	It("Resolve dependencies - cyclic dependencies", func() {
		wd, _ := os.Getwd()
		ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata"), MtaFilename: "mta_multiapps_cyclic_deps.yaml"}
		mtaStr, _ := ep.ParseFile(true)
		_, err := getModulesOrder(mtaStr)
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("eb-ui-conf-eb"))
	})

	var _ = Describe("GetModulesNames", func() {
		It("Sanity", func() {
			mtaStr := &mta.MTA{Modules: []*mta.Module{{Name: "someproj-db"}, {Name: "someproj-java"}}}
			Ω(GetModulesNames(mtaStr)).Should(Equal([]string{"someproj-db", "someproj-java"}))
		})
		It("Required module not defined", func() {
			mtaContent := readFile(getTestPath("mtahtml5", "mtaRequiredModuleNotDefined.yaml"))
			mtaStr, _ := mta.Unmarshal(mtaContent)
			_, err := GetModulesNames(mtaStr)
			Ω(err.Error()).Should(Equal(`the "abc" module is not defined`))
		})
	})
})

func readFile(file string) []byte {
	content, err := ioutil.ReadFile(file)
	Ω(err).Should(Succeed())
	s := string(content)
	s = strings.Replace(s, "\r\n", "\r", -1)
	content = []byte(s)
	return content
}

func executeAndProvideOutput(execute func() error) (string, error) {
	old := os.Stdout // keep backup of the real stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w

	err = execute()

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

	os.Stdout = old // restoring the real stdout
	// back to normal state
	_ = w.Close()
	out := <-outC
	return out, err
}

var _ = Describe("Provide", func() {
	It("Valid path to yaml", func() {

		out, err := executeAndProvideOutput(func() error {
			return ProvideModules(filepath.Join("testdata", "mtahtml5"), "dev", nil, os.Getwd)
		})
		Ω(err).Should(Succeed())
		Ω(out).Should(ContainSubstring("[ui5app ui5app2]"))
	})

	It("Invalid path to yaml", func() {
		Ω(ProvideModules(filepath.Join("testdata", "mtahtml6"), "dev", nil, os.Getwd)).Should(HaveOccurred())
	})

	It("Invalid modules dependencies", func() {
		Ω(ProvideModules(filepath.Join("testdata", "testWithWrongBuildParams"), "dev", nil, os.Getwd)).Should(HaveOccurred())
	})

	It("Invalid working folder getter", func() {
		Ω(ProvideModules("", "dev", nil, func() (string, error) {
			return "", errors.New("err")
		})).Should(HaveOccurred())
	})

})
