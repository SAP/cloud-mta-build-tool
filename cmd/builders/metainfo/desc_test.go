package metainfo

import (
	"mbtv2/cmd/mta/models"
	"os"
	"testing"
)

func Test_setManifetDesc(t *testing.T) {
	type args struct {
		file   *os.File
		mtaStr models.MTA
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setManifetDesc(tt.args.file, tt.args.mtaStr)
		})
	}
}

func TestGenMetaInf(t *testing.T) {
	type args struct {
		tmpDir string
		mtaStr models.MTA
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenMetaInf(tt.args.tmpDir, tt.args.mtaStr)
		})
	}
}
