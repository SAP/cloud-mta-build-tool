package metainfo

import (
	"bytes"
	"testing"

	"cloud-mta-build-tool/cmd/constants"
	"github.com/stretchr/testify/assert"

	"cloud-mta-build-tool/cmd/mta/models"
)

var module []string

func Test_setManifetDesc(t *testing.T) {

	tests := []struct {
		n        int
		name     string
		args     []*models.Modules
		expected []byte
	}{
		{
			n:    0,
			name: "MANIFEST.MF: One module",
			args: []*models.Modules{

				{
					Name:       "ui5",
					Type:       "html5",
					Path:       "ui5",
					Requires:   nil,
					Provides:   nil,
					Parameters: nil,
					Properties: nil,
				},
			},
			expected: []byte("Manifest-Version: 1.0\nCreated-By: SAP Application Archive Builder 0.0.1\n\n" +
				"Name: ui5" + constants.PathSep + "data.zip\nMTA-Module: ui5\nContent-Type: application/zip"),
		},
		{
			n:    0,
			name: "MANIFEST.MF: Two modules",
			args: []*models.Modules{

				{
					Name:       "ui6",
					Type:       "html5",
					Path:       "ui5",
					Requires:   nil,
					Provides:   nil,
					Parameters: nil,
					Properties: nil,
				},

				{
					Name:       "ui4",
					Type:       "html5",
					Path:       "ui5",
					Requires:   nil,
					Provides:   nil,
					Parameters: nil,
					Properties: nil,
				},
			},
			expected: []byte("Manifest-Version: 1.0\nCreated-By: SAP Application Archive Builder 0.0.1\n\n" +
				"Name: ui6" + constants.PathSep + "data.zip\nMTA-Module: ui6\nContent-Type: application/zip\n\n" +
				"Name: ui4" + constants.PathSep + "data.zip\nMTA-Module: ui4\nContent-Type: application/zip"),
		},
		{
			n:    0,
			name: "MANIFEST.MF: multi module with filter of one",
			args: []*models.Modules{

				{
					Name:       "ui6",
					Type:       "html5",
					Path:       "ui5",
					Requires:   nil,
					Provides:   nil,
					Parameters: nil,
					Properties: nil,
				},

				{
					Name:       "ui4",
					Type:       "html5",
					Path:       "ui5",
					Requires:   nil,
					Provides:   nil,
					Parameters: nil,
					Properties: nil,
				},
			},
			expected: []byte("Manifest-Version: 1.0\nCreated-By: SAP Application Archive Builder 0.0.1\n\n" +
				"Name: ui6" + constants.PathSep + "data.zip\nMTA-Module: ui6\nContent-Type: application/zip"),
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}

			// Switch was added to handle different type of slices
			switch i {
			// one module
			case 0:
				setManifetDesc(b, tt.args, module)
				if !bytes.Equal(b.Bytes(), tt.expected) {
					assert.Equal(t, string(tt.expected), b.String())
					t.Error("Test was failed")
				}
			case 1:
				//two modules
				setManifetDesc(b, tt.args, module)
				if !bytes.Equal(b.Bytes(), tt.expected) {
					assert.Equal(t, string(tt.expected), b.String())
					t.Error("Test was failed")
				}

			case 2:
				// get list of module and filter according to the name
				module = append(module, "ui6")
				setManifetDesc(b, tt.args, module)
				if !bytes.Equal(b.Bytes(), tt.expected) {
					assert.Equal(t, string(tt.expected), b.String())
					t.Error("Test was failed")
				}

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
			GenMetaInf(tt.args.tmpDir, tt.args.mtaStr, module)
		})
	}
}
