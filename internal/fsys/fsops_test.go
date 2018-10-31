package dir

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"cloud-mta-build-tool/internal/logs"

	"github.com/stretchr/testify/assert"
)

func getFullPath(relPath ...string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, filepath.Join(relPath...))
}

func TestCreateDirIfNotExist(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		dirName   string
		validator func(t *testing.T, dirName string, err error)
	}{
		{
			name:    "SanityTest",
			dirName: getFullPath("testdata", "level2", "result"),
			validator: func(t *testing.T, dirName string, err error) {
				assert.Nil(t, err)
				err = os.RemoveAll(dirName)
				assert.Nil(t, err)
			},
		},
		{
			name:    "DirectoryExists",
			dirName: getFullPath("testdata", "level2", "level3"),
			validator: func(t *testing.T, dirName string, err error) {
				assert.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateDirIfNotExist(tt.dirName)
			tt.validator(t, tt.dirName, err)
		})
	}
}

func TestArchive(t *testing.T) {
	t.Parallel()
	type args struct {
		srcFolderName  string
		archFilename   string
		archFoldername string
	}
	tests := []struct {
		name      string
		args      args
		validator func(t *testing.T, args args, err error)
	}{
		{
			name: "SanityTest",
			args: args{
				srcFolderName: getFullPath("testdata", "mtahtml5"),
				archFilename:  getFullPath("testdata", "arch.mbt"),
			},
			validator: func(t *testing.T, args args, err error) {
				assert.Nil(t, err)
				os.RemoveAll(args.archFilename)
			},
		},
		{
			name: "TargetIsNotFolder",
			args: args{
				srcFolderName: getFullPath("testdata", "level2", "level2_one.txt"),
				archFilename:  getFullPath("testdata", "arch.mbt"),
			},
			validator: func(t *testing.T, args args, err error) {
				assert.Nil(t, err)
				os.RemoveAll(args.archFilename)
			},
		},
		{
			name: "TargetIsNotExists",
			args: args{
				srcFolderName: getFullPath("testdata", "level3"),
				archFilename:  getFullPath("testdata", "arch.mbt"),
			},
			validator: func(t *testing.T, args args, err error) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Archive(tt.args.srcFolderName, tt.args.archFilename)
			tt.validator(t, tt.args, err)
		})
	}
}

func TestCreateFile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		filename  string
		validator func(t *testing.T, filename string, file *os.File, err error)
	}{
		{
			name:     "SanityTest",
			filename: "level2",
			validator: func(t *testing.T, filename string, file *os.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				file.Close()
				err = os.Remove(filename)
				assert.Nil(t, err)
			},
		},
	}

	logs.NewLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := CreateFile(tt.filename)
			tt.validator(t, tt.filename, file, err)
		})
	}
}

func countFilesInDir(name string) int {
	files, _ := ioutil.ReadDir(name)
	return len(files)
}

func TestCopyDir(t *testing.T) {

	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name         string
		args         args
		preprocessor func(t *testing.T, args args)
		validator    func(t *testing.T, args args, err error)
	}{
		{
			name:         "SanityTest",
			args:         args{getFullPath("testdata", "level2"), getFullPath("testdata", "result")},
			preprocessor: func(t *testing.T, args args) {},
			validator: func(t *testing.T, args args, err error) {
				assert.Nil(t, err)
				assert.Equal(t, countFilesInDir(args.src), countFilesInDir(args.dst))
				os.RemoveAll(args.dst)
			},
		},
		{
			name:         "SourceDirectoryDoesNotExist",
			args:         args{getFullPath("testdata", "level5"), getFullPath("testdata", "result")},
			preprocessor: func(t *testing.T, args args) {},
			validator: func(t *testing.T, args args, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name:         "SourceIsNotDirectory",
			args:         args{getFullPath("testdata", "level2", "level2_one.txt"), getFullPath("testdata", "result")},
			preprocessor: func(t *testing.T, args args) {},
			validator: func(t *testing.T, args args, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name:         "DstDirectoryNotValid",
			args:         args{getFullPath("testdata", "level2"), "/"},
			preprocessor: func(t *testing.T, args args) {},
			validator: func(t *testing.T, args args, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	logs.Logger = logs.NewLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preprocessor(t, tt.args)
			err := CopyDir(tt.args.src, tt.args.dst)
			tt.validator(t, tt.args, err)
		})
	}
}

func TestCopyDirFileLocked(t *testing.T) {
	t.Parallel()
	logs.NewLogger()
	dirName := getFullPath("testdata", "temp")
	f, _ := os.Create(dirName)

	err := CopyDir(getFullPath("testdata", "level2"), dirName)
	assert.NotNil(t, err)
	f.Close()
	os.RemoveAll(dirName)
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

func Test_copyEntries(t *testing.T) {
	t.Parallel()
	srcPath := getFullPath("testdata", "level2", "level3")
	dstPath := getFullPath("testdata", "result")
	os.MkdirAll(dstPath, os.ModePerm)
	files, _ := ioutil.ReadDir(srcPath)
	var filesWrapped [3]os.FileInfo
	for i, file := range files {
		filesWrapped[i] = testFile{file: file}
	}
	copyEntries(filesWrapped[:], srcPath, dstPath)
	assert.Equal(t, countFilesInDir(srcPath)-1, countFilesInDir(dstPath))
	os.RemoveAll(dstPath)

	dstPath = getFullPath("testdata", "//")
	err := copyEntries(filesWrapped[:], getFullPath("testdata", "level2", "levelx"), dstPath)
	assert.NotNil(t, err)
}

func Test_copyFile(t *testing.T) {
	t.Parallel()
	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name      string
		args      args
		validator func(t *testing.T, args args, err error)
	}{
		{
			name: "SourceNotExists",
			args: args{getFullPath("testdata", "fileSrc"), getFullPath("testdata", "fileDst")},
			validator: func(t *testing.T, args args, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "SourceIsDirectory",
			args: args{getFullPath("testdata", "level2"), getFullPath("testdata", "level2", "fileDst")},
			validator: func(t *testing.T, args args, err error) {
				assert.NotNil(t, err)
				os.RemoveAll(args.dst)
			},
		},
		{
			name: "WrongDestinationName",
			args: args{getFullPath("testdata", "level2", "level2_one.txt"), getFullPath("testdata", "level2", "/")},
			validator: func(t *testing.T, args args, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "DestinationExists",
			args: args{getFullPath("testdata", "level2", "level3", "level3_one.txt"), getFullPath("testdata", "level2", "level3", "level3_two.txt")},
			validator: func(t *testing.T, args args, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := copyFile(tt.args.src, tt.args.dst)
			tt.validator(t, tt.args, err)
		})
	}
}
