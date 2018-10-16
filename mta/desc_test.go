package mta

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"cloud-mta-build-tool/cmd/fsys"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var module []string

func Test_setManifetDesc(t *testing.T) {

	tests := []struct {
		n        int
		name     string
		args     []*Modules
		expected []byte
	}{
		{
			n:    0,
			name: "MANIFEST.MF: One module",
			args: []*Modules{

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
			expected: []byte("manifest-Version: 1.0\nCreated-By: SAP Application Archive Builder 0.0.1\n\n" +
				"Name: ui5/data.zip\nMTA-Module: ui5\nContent-Type: application/zip"),
		},
		{
			n:    0,
			name: "MANIFEST.MF: Two modules",
			args: []*Modules{

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
			expected: []byte("manifest-Version: 1.0\nCreated-By: SAP Application Archive Builder 0.0.1\n\n" +
				"Name: ui6/data.zip\nMTA-Module: ui6\nContent-Type: application/zip\n\n" +
				"Name: ui4/data.zip\nMTA-Module: ui4\nContent-Type: application/zip"),
		},
		{
			n:    0,
			name: "MANIFEST.MF: multi module with filter of one",
			args: []*Modules{

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
			expected: []byte("manifest-Version: 1.0\nCreated-By: SAP Application Archive Builder 0.0.1\n\n" +
				"Name: ui6/data.zip\nMTA-Module: ui6\nContent-Type: application/zip"),
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
				// two modules
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
	var mtaSingleModule = []byte(`
_schema-version: "2.0.0"
ID: com.sap.webide.feature.management
version: 1.0.0

modules:
  - name: htmlapp
    type: html5
    path: app
`)
	mta := MTA{}

	wd, _ := os.Getwd()

	yaml.Unmarshal(mtaSingleModule, &mta)
	GenMetaInfo(filepath.Join(wd, "testdata"), mta, []string{"htmlapp"}, func(mtaStr MTA) {

	})

	_, err := ioutil.ReadFile(filepath.Join(wd, "testdata", "META-INF", "MANIFEST.MF"))
	assert.Nil(t, err)
	_, err = ioutil.ReadFile(filepath.Join(wd, "testdata", "META-INF", "mtad.yaml"))
	assert.Nil(t, err)
	os.RemoveAll(filepath.Join(dir.GetPath(), "testdata", "META-INF"))
}
