package dir

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/types"
)

func getFullPath(relPath ...string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, filepath.Join(relPath...))
}

type testMtaYamlStr struct {
	fullpath string
	path     string
	err      error
}

func (t *testMtaYamlStr) GetMtaYamlFilename() string {
	return t.fullpath
}

func (t *testMtaYamlStr) GetMtaYamlPath() string {
	return t.path
}

func (t *testMtaYamlStr) GetMtaExtYamlPath(platform string) string {
	return t.fullpath
}

var _ = Describe("FSOPS", func() {

	var _ = Describe("CreateDir", func() {

		AfterEach(func() {
			os.RemoveAll(getFullPath("testdata", "level2", "result"))
		})

		var _ = DescribeTable("CreateDir", func(dirPath string) {

			Ω(CreateDirIfNotExist(dirPath)).Should(Succeed())
		},
			Entry("Sanity", getFullPath("testdata", "level2", "result")),
			Entry("DirectoryExists", getFullPath("testdata", "level2", "level3")),
		)
		It("Fails because file with the same name exists", func() {
			CreateDirIfNotExist(getFullPath("testdata", "level2", "result"))
			file, _ := os.Create(getFullPath("testdata", "level2", "result", "file"))
			file.Close()
			Ω(CreateDirIfNotExist(getFullPath("testdata", "level2", "result", "file"))).Should(HaveOccurred())
		})
	})

	var _ = Describe("Archive", func() {
		var targetFilePath = getFullPath("testdata", "arch.mbt")

		AfterEach(func() {
			os.RemoveAll(targetFilePath)
		})

		var _ = DescribeTable("Archive", func(source, target string, matcher GomegaMatcher, created bool) {

			Ω(Archive(source, target, nil)).Should(matcher)
			if created {
				Ω(target).Should(BeAnExistingFile())
			} else {
				Ω(target).ShouldNot(BeAnExistingFile())
			}
		},
			Entry("Sanity",
				getFullPath("testdata", "mtahtml5"), targetFilePath, Succeed(), true),
			Entry("SourceIsNotFolder",
				getFullPath("testdata", "level2", "level2_one.txt"), targetFilePath, Succeed(), true),
			Entry("Target is empty string",
				getFullPath("testdata", "mtahtml5"), "", HaveOccurred(), false),
			Entry("Source is empty string",
				"", "", HaveOccurred(), false),
		)
	})

	var _ = Describe("Create File", func() {
		AfterEach(func() {
			os.RemoveAll(getFullPath("testdata", "result.txt"))
		})
		It("Sanity", func() {
			file, err := CreateFile(getFullPath("testdata", "result.txt"))
			Ω(getFullPath("testdata", "result.txt")).Should(BeAnExistingFile())
			file.Close()
			Ω(err).Should(BeNil())
		})
		It("Fails on empty path", func() {
			_, err := CreateFile("")
			Ω(err).Should(HaveOccurred())
		})
	})

	var _ = Describe("CopyDir", func() {
		var targetPath = getFullPath("testdata", "result")
		AfterEach(func() {
			os.RemoveAll(targetPath)
		})

		It("Sanity - parallel", func() {
			sourcePath := getFullPath("testdata", "level2")
			Ω(CopyDir(sourcePath, targetPath, true, CopyEntriesInParallel)).Should(Succeed())
			Ω(countFilesInDir(targetPath)).Should(Equal(countFilesInDir(sourcePath)))
		})

		It("Sanity - not parallel", func() {
			sourcePath := getFullPath("testdata", "level2")
			Ω(CopyDir(sourcePath, targetPath, true, CopyEntries)).Should(Succeed())
			Ω(countFilesInDir(targetPath)).Should(Equal(countFilesInDir(sourcePath)))
		})

		It("TargetFileLocked", func() {
			f, _ := os.Create(targetPath)
			sourcePath := getFullPath("testdata", "level2")
			Ω(CopyDir(sourcePath, targetPath, true, CopyEntries)).Should(HaveOccurred())
			f.Close()
		})

		It("TargetFileLocked", func() {
			f, _ := os.Create(targetPath)
			sourcePath := getFullPath("testdata", "level2")
			Ω(CopyDir(sourcePath, targetPath, true, CopyEntriesInParallel)).Should(HaveOccurred())
			f.Close()
		})

		var _ = DescribeTable("Invalid cases", func(source, target string) {
			Ω(CopyDir(source, targetPath, true, CopyEntries)).Should(HaveOccurred())
		},
			Entry("SourceDirectoryDoesNotExist", getFullPath("testdata", "level5"), targetPath),
			Entry("SourceIsNotDirectory", getFullPath("testdata", "level2", "level2_one.txt"), targetPath),
			Entry("DstDirectoryNotValid", getFullPath("level2"), ":"),
		)

		var _ = DescribeTable("Invalid cases - parallel", func(source, target string) {
			Ω(CopyDir(source, targetPath, true, CopyEntriesInParallel)).Should(HaveOccurred())
		},
			Entry("SourceDirectoryDoesNotExist", getFullPath("testdata", "level5"), targetPath),
			Entry("SourceIsNotDirectory", getFullPath("testdata", "level2", "level2_one.txt"), targetPath),
			Entry("DstDirectoryNotValid", getFullPath("level2"), ":"),
		)

		var _ = DescribeTable("Copy File - Invalid", func(source, target string, matcher GomegaMatcher) {
			Ω(CopyFile(source, target)).Should(matcher)
		},
			Entry("SourceNotExists", getFullPath("testdata", "fileSrc"), targetPath, HaveOccurred()),
			Entry("SourceIsDirectory", getFullPath("testdata", "level2"), targetPath, HaveOccurred()),
			Entry("WrongDestinationName", getFullPath("testdata", "level2", "level2_one.txt"), getFullPath("testdata", "level2", "/"), HaveOccurred()),
			Entry("DestinationExists", getFullPath("testdata", "level2", "level3", "level3_one.txt"), getFullPath("testdata", "level2", "level3", "level3_two.txt"), Succeed()),
		)
		var _ = DescribeTable("Copy File - Invalid", func(source, target string, matcher GomegaMatcher) {
			Ω(CopyFileWithMode(source, target, os.ModePerm)).Should(matcher)
		},
			Entry("TargetPathEmpty", getFullPath("testdata", "fileSrc"), "", HaveOccurred()),
			Entry("SourceIsDirectory", getFullPath("testdata", "level2"), targetPath, HaveOccurred()),
			Entry("DestinationExists", getFullPath("testdata", "level2", "level3", "level3_one.txt"), getFullPath("testdata", "level2", "level3", "level3_two.txt"), Succeed()),
		)
	})

	var _ = Describe("Copy Entries", func() {

		AfterEach(func() {
			os.RemoveAll(getFullPath("testdata", "result"))
		})

		It("Sanity", func() {
			sourcePath := getFullPath("testdata", "level2", "level3")
			targetPath := getFullPath("testdata", "result")
			os.MkdirAll(targetPath, os.ModePerm)
			files, _ := ioutil.ReadDir(sourcePath)
			// Files wrapped to overwrite their methods
			var filesWrapped []os.FileInfo
			Ω(CopyEntries(filesWrapped, sourcePath, targetPath)).Should(Succeed())
			for _, file := range files {
				filesWrapped = append(filesWrapped, testFile{file: file})
			}
			Ω(CopyEntries(filesWrapped, sourcePath, targetPath)).Should(Succeed())
			Ω(countFilesInDir(sourcePath) - 1).Should(Equal(countFilesInDir(targetPath)))
			os.RemoveAll(targetPath)
		})
		It("Sanity - copy in parallel", func() {
			sourcePath := getFullPath("testdata", "level2", "level3")
			targetPath := getFullPath("testdata", "result")
			os.MkdirAll(targetPath, os.ModePerm)
			files, _ := ioutil.ReadDir(sourcePath)
			// Files wrapped to overwrite their methods
			var filesWrapped []os.FileInfo
			Ω(CopyEntriesInParallel(filesWrapped, sourcePath, targetPath)).Should(Succeed())
			for _, file := range files {
				filesWrapped = append(filesWrapped, testFile{file: file})
			}
			Ω(CopyEntriesInParallel(filesWrapped, sourcePath, targetPath)).Should(Succeed())
			Ω(countFilesInDir(sourcePath) - 1).Should(Equal(countFilesInDir(targetPath)))
			os.RemoveAll(targetPath)
		})
	})

	var _ = Describe("Copy By Patterns", func() {

		AfterEach(func() {
			os.RemoveAll(getFullPath("testdata", "result"))
		})

		var _ = DescribeTable("Valid Cases", func(modulePath string, patterns, expectedFiles []string) {
			sourcePath := getFullPath("testdata", "testbuildparams", modulePath)
			targetPath := getFullPath("testdata", "result")
			Ω(CopyByPatterns(sourcePath, targetPath, patterns)).Should(Succeed())
			for _, file := range expectedFiles {
				Ω(file).Should(BeAnExistingFile())
			}
		},
			Entry("Single file", "ui2",
				[]string{"deep/folder/inui2/anotherfile.txt"},
				[]string{getFullPath("testdata", "result", "anotherfile.txt")}),
			Entry("Wildcard for 2 files", "ui2",
				[]string{"deep/*/inui2/another*"},
				[]string{getFullPath("testdata", "result", "anotherfile.txt"),
					getFullPath("testdata", "result", "anotherfile2.txt")}),
			Entry("Wildcard for 2 files - dot start", "ui2",
				[]string{"./deep/*/inui2/another*"},
				[]string{getFullPath("testdata", "result", "anotherfile.txt"),
					getFullPath("testdata", "result", "anotherfile2.txt")}),
			Entry("Specific folder of second level", "ui2",
				[]string{"*/folder/*"},
				[]string{
					getFullPath("testdata", "result", "inui2", "anotherfile.txt"),
					getFullPath("testdata", "result", "inui2", "anotherfile2.txt")}),
			Entry("All", "ui1",
				[]string{"*"},
				[]string{getFullPath("testdata", "result", "webapp", "Component.js")}),
			Entry("Dot", "ui1",
				[]string{"."},
				[]string{getFullPath("testdata", "result", "ui1", "webapp", "Component.js")}),
			Entry("Multiple patterns", "ui2", //
				[]string{"deep/folder/inui2/anotherfile.txt", "*/folder/"},
				[]string{
					getFullPath("testdata", "result", "folder", "inui2", "anotherfile.txt"),
					getFullPath("testdata", "result", "anotherfile.txt")}),
			Entry("Empty patterns", "ui2",
				[]string{},
				[]string{}),
		)

		var _ = DescribeTable("Invalid Cases", func(targetPath, modulePath string, patterns []string) {
			sourcePath := getFullPath("testdata", "testbuildparams", modulePath)
			err := CopyByPatterns(sourcePath, targetPath, patterns)
			Ω(err).Should(HaveOccurred())
		},
			Entry("Target path relates to file ",
				getFullPath("testdata", "testbuildparams", "mta.yaml"), "ui2",
				[]string{"deep/folder/inui2/somefile.txt"}),
			Entry("Wrong pattern ",
				getFullPath("testdata", "result"), "ui2", []string{"[a,b"}),
			Entry("Empty target path", "", "ui2", []string{"[a,b"}),
		)
	})

	It("getRelativePath", func() {
		Ω(getRelativePath(getFullPath("abc", "xyz", "fff"),
			filepath.Join(getFullPath()))).Should(Equal(string(filepath.Separator) + filepath.Join("abc", "xyz", "fff")))
	})

	It("copyByPattern - fails because target is file", func() {
		Ω(copyByPattern(getPath("testdata", "mtahtml5"),
			getFullPath("testdata", "level2", "level2_one.txt"), "m*")).Should(HaveOccurred())
	})

	It("copyEntries - fails on entry with empty name", func() {
		Ω(copyEntries([]string{""}, getPath("testdata", "mtahtml5"),
			getFullPath("testdata", "level2", "level2_one.txt"), "m*")).Should(HaveOccurred())
	})

	It("changeTargetMode - fails if source does not exist", func() {
		Ω(changeTargetMode(getPath("testdata", "not-exists"), getPath("testdata", "not-exists-2"))).
			Should(HaveOccurred())
	})

	var _ = Describe("Read", func() {
		It("Sanity", func() {
			test := testMtaYamlStr{
				fullpath: getFullPath("testdata", "testproject", "mta.yaml"),
				path:     getFullPath("testdata", "testproject", "mta.yaml"),
				err:      nil,
			}
			res, resErr := Read(&test)
			Ω(res).ShouldNot(BeNil())
			Ω(resErr).Should(BeNil())
		})
	})

	var _ = Describe("ReadExt", func() {
		It("Sanity", func() {
			test := testMtaYamlStr{
				fullpath: getFullPath("testdata", "testproject", "mta.yaml"),
				path:     getFullPath("testdata", "testproject", "mta.yaml"),
				err:      nil,
			}
			res, resErr := ReadExt(&test, "cf")
			Ω(res).ShouldNot(BeNil())
			Ω(resErr).Should(BeNil())
		})
	})

	var _ = DescribeTable("CloseFile", func(toFail bool, errorArg error, expectedErr error) {
		testFile := testCloser{fail: toFail}
		if expectedErr == nil {
			Ω(CloseFile(&testFile, errorArg)).Should(BeNil())
		} else {
			Ω(CloseFile(&testFile, errorArg).Error()).Should(Equal(expectedErr.Error()))
		}
	},
		Entry("No error", false, nil, nil),
		Entry("Original error only", false, errors.New("original error"), errors.New("original error")),
		Entry("New error", true, nil, errors.New("failed to close")),
		Entry("Original and new errors", true, errors.New("original error"), errors.New("original error")),
	)
})

func countFilesInDir(name string) int {
	files, _ := ioutil.ReadDir(name)
	return len(files)
}

type testFile struct {
	file os.FileInfo
}

func (file testFile) Name() string {
	return file.file.Name()
}

func (file testFile) Size() int64 {
	return file.file.Size()
}

func (file testFile) Mode() os.FileMode {
	if strings.Contains(file.file.Name(), "level3_one.txt") {
		return os.ModeSymlink
	}
	return file.file.Mode()
}

func (file testFile) ModTime() time.Time {
	return file.file.ModTime()
}

func (file testFile) IsDir() bool {
	return file.file.IsDir()
}

func (file testFile) Sys() interface{} {
	return nil
}

type testCloser struct {
	fail bool
}

func (f *testCloser) Close() error {
	if f.fail {
		return errors.New("failed to close")
	}
	return nil
}
