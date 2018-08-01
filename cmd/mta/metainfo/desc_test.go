package metainfo

import (
	"mbtv2/cmd/mta/models"
	//"os"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setManifetDesc(t *testing.T) {

	tests := []struct {
		n        int
		name     string
		args     models.MTA
		expected []byte
	}{
		{
			name: "name",

			//	TODO - read the version via semver
			expected: []byte("Manifest-Version: 1.0 \nCreated-By: SAP Application Archive Builder 0.0.1"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			setManifetDesc(b, tt.args)
			if !bytes.Equal(b.Bytes(), tt.expected) {
				assert.Equal(t, string(tt.expected), b.String())

				//fmt.Printf("expect:\n%v\nactual:%v\n\n", b.String(), string(tt.expected))
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
