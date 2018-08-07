package dir

import "testing"

func TestGetPath(t *testing.T) {
	tests := []struct {
		name    string
		wantDir string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDir := GetPath(); gotDir != tt.wantDir {
				t.Errorf("GetPath() = %v, want %v", gotDir, tt.wantDir)
			}
		})
	}
}

func TestProjectPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProjectPath(); got != tt.want {
				t.Errorf("ProjectPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
