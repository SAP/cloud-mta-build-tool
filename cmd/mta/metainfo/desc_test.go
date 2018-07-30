package metainfo

import (
	"mbtv2/cmd/mta/models"
	"os"
	"testing"
	"bytes"
)

func Test_setManifetDesc(t *testing.T) {

	type args struct {
		file   *os.File
		mtaStr models.MTA
	}
	var tests []struct {
		name     string
		args     args
		expected []byte
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			setManifetDesc(tt.args.file, tt.args.mtaStr)
			if !bytes.Equal(b.Bytes(), tt.expected) {
				t.Error("Fail")
			}

		})
	}
}

func TestGenMetaInf(t *testing.T) {
	type args struct {
		tmpDir string
		mtaStr models.MTA
	}
	var tests []struct {
		name string
		args args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenMetaInf(tt.args.tmpDir, tt.args.mtaStr)
		})
	}
}
