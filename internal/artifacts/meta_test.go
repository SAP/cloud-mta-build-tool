package artifacts

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/platform"
	"github.com/SAP/cloud-mta-build-tool/internal/version"
	"github.com/SAP/cloud-mta/mta"
)

var _ = Describe("Meta", func() {

	AfterEach(func() {
		Ω(os.RemoveAll(getResultPath())).Should(Succeed())
	})

	var _ = Describe("ExecuteGenMeta", func() {

		It("Sanity", func() {
			createMtahtml5TmpFolder()
			Ω(ExecuteGenMeta(getTestPath("mtahtml5"), getResultPath(), "dev", nil, "CF", os.Getwd)).Should(Succeed())
			Ω(getFullPathInTmpFolder("mtahtml5", "META-INF", "MANIFEST.MF")).Should(BeAnExistingFile())
			Ω(getFullPathInTmpFolder("mtahtml5", "META-INF", "mtad.yaml")).Should(BeAnExistingFile())
		})

		It("Fails on META-INF folder creation", func() {
			createDirInTmpFolder("mtahtml5")
			createFileInTmpFolder("mtahtml5", "META-INF")
			err := ExecuteGenMeta(getTestPath("mtahtml5"), getResultPath(), "dev", nil, "CF", os.Getwd)
			checkError(err, dir.FolderCreationFailedMsg, getFullPathInTmpFolder("mtahtml5", "META-INF"))
		})

		It("Wrong location - fails on Working directory get", func() {
			err := ExecuteGenMeta("", "", "dev", nil, "cf", func() (string, error) {
				return "", errors.New("error of working dir get")
			})
			checkError(err, "error of working dir get")
		})
		It("Wrong platform", func() {
			createDirInTmpFolder("mtahtml5", "ui5app2")
			createDirInTmpFolder("mtahtml5", "testapp")
			err := ExecuteGenMeta(getTestPath("mtahtml5"), getResultPath(), "dev", nil, "xx", os.Getwd)
			checkError(err, invalidPlatformMsg, "xx")
		})
		It("generateMeta fails on wrong source path - parse mta fails", func() {
			err := ExecuteGenMeta(getTestPath("mtahtml6"), getResultPath(), "dev", nil, "cf", os.Getwd)
			checkError(err, getTestPath("mtahtml6", "mta.yaml"))
		})
	})

	var _ = Describe("GenMetaInfo", func() {
		ep := dir.Loc{SourcePath: getTestPath("testproject"), TargetPath: getResultPath()}
		var mtaSingleModule = []byte(`
_schema-version: "2.0.0"
ID: mta_proj
version: 1.0.0

modules:
  - name: htmlapp
    type: html5
    path: app
`)

		It("Sanity", func() {
			m, err := mta.Unmarshal(mtaSingleModule)
			Ω(err).Should(Succeed())
			createDirInTmpFolder("testproject", "htmlapp")
			createFileInTmpFolder("testproject", "htmlapp", "data.zip")
			createDirInTmpFolder("testproject", "META-INF")
			Ω(genMetaInfo(&ep, &ep, &ep, ep.IsDeploymentDescriptor(), "cf", m, true, true)).Should(Succeed())
			Ω(ep.GetManifestPath()).Should(BeAnExistingFile())
			Ω(ep.GetMtadPath()).Should(BeAnExistingFile())
		})

		It("Meta creation fails - fails on conversion by platform", func() {
			m, err := mta.Unmarshal(mtaSingleModule)
			Ω(err).Should(Succeed())
			createDirInTmpFolder("testproject", "app")
			createDirInTmpFolder("testproject", "META-INF")
			cfg := platform.PlatformConfig
			platform.PlatformConfig = []byte(`very bad config`)
			Ω(genMetaInfo(&ep, &ep, &ep, ep.IsDeploymentDescriptor(), "cf", m, true, true)).Should(HaveOccurred())
			platform.PlatformConfig = cfg
		})

		It("Fails on create file for manifest path", func() {
			loc := testLoc{ep}
			m, err := mta.Unmarshal(mtaSingleModule)
			Ω(err).Should(Succeed())
			Ω(genMetaInfo(&loc, &ep, &ep, ep.IsDeploymentDescriptor(), "cf", m, true, true)).Should(HaveOccurred())
		})

		var _ = Describe("Fails on setManifestDesc", func() {
			var config []byte

			BeforeEach(func() {
				config = make([]byte, len(version.VersionConfig))
				copy(config, version.VersionConfig)
				// Simplified commands configuration (performance purposes). removed "npm prune --production"
				version.VersionConfig = []byte(`
cli_version:["x"]
`)
			})

			AfterEach(func() {
				version.VersionConfig = make([]byte, len(config))
				copy(version.VersionConfig, config)
				Ω(os.RemoveAll(getResultPath())).Should(Succeed())
			})

			It("Fails on get version", func() {
				m, err := mta.Unmarshal(mtaSingleModule)
				Ω(err).Should(Succeed())
				Ω(genMetaInfo(&ep, &ep, &ep, ep.IsDeploymentDescriptor(), "cf", m, true, true)).Should(HaveOccurred())
			})
		})
	})

	var _ = Describe("Generate Commands", func() {

		AfterEach(func() {
			Ω(os.RemoveAll(getResultPath())).Should(Succeed())
		})

		readFileContent := func(ep dir.IMtaParser) *mta.MTA {
			mtaObj, _ := ep.ParseFile()
			return mtaObj
		}

		It("Generate Meta", func() {
			createMtahtml5TmpFolder()
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			Ω(generateMeta(&ep, &ep, false, "cf", true, true)).Should(Succeed())
			Ω(readFileContent(&dir.Loc{SourcePath: getFullPathInTmpFolder("mtahtml5", "META-INF"), Descriptor: "dep"})).
				Should(Equal(readFileContent(&dir.Loc{SourcePath: getTestPath("golden"), Descriptor: "dep"})))
		})

		It("Generate Meta - fails on missing module path in temporary folder", func() {
			createMtahtml5WithMissingModuleTmpFolder()
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			Ω(generateMeta(&ep, &ep, false, "cf", false, true)).Should(HaveOccurred())
		})

		It("Generate Meta - doesn't fail on missing module path in temporary folder because module configured as no-source", func() {
			createMtahtml5WithMissingModuleTmpFolder()
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath(), MtaFilename: "mtaWithNoSource.yaml"}
			Ω(generateMeta(&ep, &ep, false, "cf", true, true)).Should(Succeed())
			Ω(readFileContent(&dir.Loc{SourcePath: getFullPathInTmpFolder("mtahtml5", "META-INF"), Descriptor: "dep"})).
				Should(Equal(readFileContent(&dir.Loc{SourcePath: getTestPath("goldenNoSource"), Descriptor: "dep"})))
		})

		It("Generate Meta - mta not exists", func() {
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath(),
				MtaFilename: "mtaNotExists.yaml"}
			err := generateMeta(&ep, &ep, false, "cf", true, true)
			checkError(err, ep.GetMtaYamlPath())
		})

		Describe("mocking platform", func() {

			var platformConfig []byte

			BeforeEach(func() {
				platformConfig = platform.PlatformConfig
				platform.PlatformConfig = []byte("wrong config")
			})

			AfterEach(func() {
				platform.PlatformConfig = platformConfig
			})

			It("Generate Meta fails on platform parsing", func() {
				createMtahtml5TmpFolder()
				ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
				err := generateMeta(&ep, &ep, false, "cf", true, true)
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(ContainSubstring(fmt.Sprintf(genMTADTypeTypeCnvMsg, "cf")))
				Ω(err.Error()).Should(ContainSubstring(platform.UnmarshalFailedMsg))
			})
		})

		It("Generate Mtar", func() {
			createMtahtml5TmpFolder()
			ep := dir.Loc{SourcePath: getTestPath("mtahtml5"), TargetPath: getResultPath()}
			err := generateMeta(&ep, &ep, false, "cf", true, true)
			Ω(err).Should(Succeed())
			mtarPath, err := generateMtar(&ep, &ep, &ep, true, "")
			Ω(err).Should(Succeed())
			Ω(mtarPath).Should(BeAnExistingFile())
		})
	})
	Describe("ExecuteMerge", func() {
		resultFileName := "result.yaml"
		resultFilePath := getTestPath("result", resultFileName)

		It("Succeeds with single mtaext file", func() {
			err := ExecuteMerge(getTestPath("mta_with_ext"), getResultPath(), []string{"cf-mtaext.mtaext"}, resultFileName, os.Getwd)
			Ω(err).Should(Succeed())
			Ω(resultFilePath).Should(BeAnExistingFile())
			compareMTAContent(getTestPath("mta_with_ext", "golden1.yaml"), resultFilePath)
		})
		It("Succeeds with two mtaext files", func() {
			err := ExecuteMerge(getTestPath("mta_with_ext"), getResultPath(), []string{"other.mtaext", "cf-mtaext.mtaext"}, resultFileName, os.Getwd)
			Ω(err).Should(Succeed())
			Ω(resultFilePath).Should(BeAnExistingFile())
			compareMTAContent(getTestPath("mta_with_ext", "golden2.yaml"), resultFilePath)
		})
		It("Fails when the result file name is not sent", func() {
			err := ExecuteMerge(getTestPath("mta_with_ext"), getResultPath(), []string{"cf-mtaext.mtaext"}, "", os.Getwd)
			checkError(err, mergeNameRequiredMsg)
		})
		It("Fails when the result file already exists", func() {
			Ω(dir.CreateDirIfNotExist(getResultPath())).Should(Succeed())
			createFileInGivenPath(resultFilePath)

			err := ExecuteMerge(getTestPath("mta_with_ext"), getResultPath(), []string{"cf-mtaext.mtaext"}, resultFileName, os.Getwd)
			checkError(err, mergeFailedOnFileCreationMsg, resultFilePath)
		})
		It("Fails when the result directory is a file", func() {
			Ω(dir.CreateDirIfNotExist(filepath.Dir(getResultPath()))).Should(Succeed())
			createFileInGivenPath(getResultPath())

			err := ExecuteMerge(getTestPath("mta_with_ext"), getResultPath(), []string{"cf-mtaext.mtaext"}, resultFileName, os.Getwd)
			checkError(err, dir.FolderCreationFailedMsg, getResultPath())
		})
		It("Fails when the mtaext file doesn't exist", func() {
			err := ExecuteMerge(getTestPath("mta_with_ext"), getResultPath(), []string{"invalid.yaml"}, resultFileName, os.Getwd)
			checkError(err, getTestPath("mta_with_ext", "invalid.yaml"))
		})
		It("Fails when wdGetter fails", func() {
			err := ExecuteMerge("", getResultPath(), []string{"cf-mtaext.yaml"}, resultFileName, func() (string, error) {
				return "", errors.New("an error occurred")
			})
			checkError(err, "an error occurred")
		})
	})
})

type testLoc struct {
	loc dir.Loc
}

func (loc *testLoc) GetMetaPath() string {
	return loc.loc.GetMetaPath()
}

func (loc *testLoc) GetMtadPath() string {
	return loc.loc.GetMtadPath()
}

func (loc *testLoc) GetManifestPath() string {
	return filepath.Join(loc.loc.GetManifestPath(), "folderNotExists", "MANIFEST.MF")
}

func (loc *testLoc) GetMtarDir(targetProvided bool) string {
	return loc.loc.GetMtarDir(targetProvided)
}

func (loc *testLoc) GetSourceModuleDir(modulePath string) string {
	return loc.loc.GetSourceModuleDir(modulePath)
}

func (loc *testLoc) GetTargetModuleDir(moduleName string) string {
	return loc.loc.GetTargetModuleDir(moduleName)
}

func (loc *testLoc) GetSourceModuleArtifactRelPath(modulePath, artifactPath string) (string, error) {
	return loc.loc.GetSourceModuleArtifactRelPath(modulePath, artifactPath)
}

func (loc *testLoc) GetTargetTmpDir() string {
	return loc.loc.GetTargetTmpDir()
}

func (loc *testLoc) GetTargetTmpRoot() string {
	return loc.loc.GetTargetTmpRoot()
}

func createMtahtml5TmpFolder() {
	createDirInTmpFolder("mtahtml5", "ui5app2")
	createDirInTmpFolder("mtahtml5", "ui5app")
	createDirInTmpFolder("mtahtml5", "META-INF")
	createFileInTmpFolder("mtahtml5", "ui5app", "data.zip")
	createFileInTmpFolder("mtahtml5", "ui5app2", "data.zip")
	createFileInTmpFolder("mtahtml5", "xs-security.json")
}

func createMtahtml5WithMissingModuleTmpFolder() {
	createDirInTmpFolder("mtahtml5", "ui5app2")
	createDirInTmpFolder("mtahtml5", "META-INF")
	createFileInTmpFolder("mtahtml5", "ui5app2", "data.zip")
	createFileInTmpFolder("mtahtml5", "xs-security.json")
}

func compareMTAContent(expectedFileName string, actualFileName string) {
	actual, err := ioutil.ReadFile(expectedFileName)
	Ω(err).Should(Succeed())
	actualMta, err := mta.Unmarshal(actual)
	Ω(err).Should(Succeed())
	expected, err := ioutil.ReadFile(actualFileName)
	Ω(err).Should(Succeed())
	expectedMta, err := mta.Unmarshal(expected)
	Ω(err).Should(Succeed())
	Ω(actualMta).Should(Equal(expectedMta))
}
