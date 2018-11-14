package proc

import (
	"reflect"
	"runtime"
	"testing"
)

func TestOsCore(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want []proc
	}{
		{
			name: "Operating system core",
			want: getCores(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OsCore()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OsCore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getCores() []proc {

	switch runtime.GOOS {
	case "linux":
		return []proc{{`NPROCS = $(shell grep -c 'processor' /proc/cpuinfo)`, `MAKEFLAGS += -j$(NPROCS)`}}
	case "darwin":
		return []proc{{`NPROCS = $(sysctl -n hw.ncpu)`, `MAKEFLAGS += -j$(NPROCS)`}}
	case "windows":
		return []proc{{`NPROCS = $(shell echo %NUMBER_OF_PROCESSORS%)`, `MAKEFLAGS += -j$(NPROCS)`}}
	default:
		return []proc{}
	}
}
