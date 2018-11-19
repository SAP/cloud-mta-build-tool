package proc

import "runtime"

type Proc struct {
	NPROCS    string
	MAKEFLAGS string
}

// OsCore - Get available cores according to the running OS
func OsCore() []Proc {
	switch runtime.GOOS {
	case "linux":
		return []Proc{{`NPROCS = $(shell grep -c 'processor' /proc/cpuinfo)`, `MAKEFLAGS += -j$(NPROCS)`}}
	case "darwin":
		return []Proc{{`NPROCS = $(sysctl -n hw.ncpu)`, `MAKEFLAGS += -j$(NPROCS)`}}
	case "windows":
		return []Proc{{`NPROCS = $(shell echo %NUMBER_OF_PROCESSORS%)`, `MAKEFLAGS += -j$(NPROCS)`}}
	default:
		return []Proc{}
	}
}
