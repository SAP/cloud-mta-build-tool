package mta

import (
	"testing"

	"cloud-mta-build-tool/cmd/mta/models"
)

func TestSetMtaProp(t *testing.T) {
	type args struct {
		mtaStruct models.MTA
	}
	var tests []struct {
		name string
		args args
		want string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetMtaProp(tt.args.mtaStruct); got != tt.want {
				t.Errorf("SetMtaProp() = %v, want %v", got, tt.want)
			}
		})
	}
}
