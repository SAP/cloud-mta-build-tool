package mta

import (
	"cloud-mta-build-tool/cmd/mta/models"
	"testing"
)

func TestSetMtaProp(t *testing.T) {
	type args struct {
		mtaStruct models.MTA
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetMtaProp(tt.args.mtaStruct); got != tt.want {
				t.Errorf("SetMtaProp() = %v, want %v", got, tt.want)
			}
		})
	}
}
