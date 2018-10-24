package dir

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Path", func() {
	It("GetPath", func() {
		Ω(GetCurrentPath()).Should(BeADirectory())
	})

	It("GetFullPath Method", func() {
		currentPath, _ := GetCurrentPath()
		Ω(Path{currentPath}.GetFullPath("testdata", "mtahtml5")).Should(BeADirectory())
		Ω(Path{currentPath}.GetFullPath("testdata", "level2", "level2_one.txt")).Should(BeAnExistingFile())
	})

	It("GetFullPath Function", func() {
		Ω(GetFullPath("testdata", "mtahtml5")).Should(BeADirectory())
		Ω(GetFullPath("testdata", "level2", "level2_one.txt")).Should(BeAnExistingFile())
	})

	It("GetArtifactsPath", func() {
		Ω(GetArtifactsPath(getFullPath())).ShouldNot(BeADirectory())
		Ω(GetArtifactsPath(getFullPath("testdata", "level2", "level3"))).Should(BeADirectory())
	})

	It("GetRelativePath", func() {
		Ω(GetRelativePath(getFullPath("abc", "xyz", "fff"),
			filepath.Join(getFullPath()))).Should(Equal(string(filepath.Separator) + filepath.Join("abc", "xyz", "fff")))
	})
})
