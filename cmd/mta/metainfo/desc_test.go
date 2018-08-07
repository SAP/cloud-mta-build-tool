package metainfo

import (
	//"os"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"cloud-mta-build-tool/cmd/mta/models"
)

func Test_setManifetDesc(t *testing.T) {

	tests := []struct {
		n        int
		name     string
		args     []*models.Modules
		expected []byte
	}{
		{
			n:    0,
			name: "name",
			args: []*models.Modules{

				{
					Name:        "ui5",
					Type:        "html5",
					Path:        "ui5",
					Requires:    nil,
					Provides:    nil,
					Parameters:  nil,
					BuildParams: nil,
					Properties:  nil,
				},
			},
			expected: []byte("Manifest-Version: 1.0 \nCreated-By: SAP Application Archive Builder 0.0.1\n\nName: ui5/data.zip\nMTA-Module: ui5\nContent-Type: application/zip"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			setManifetDesc(b, tt.args)
			if !bytes.Equal(b.Bytes(), tt.expected) {
				assert.Equal(t, string(tt.expected), b.String())
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
