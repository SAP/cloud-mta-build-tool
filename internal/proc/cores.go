package proc

import "runtime"

type proc struct {
	NPROCS    string
	MAKEFLAGS string
}

// OsCore - Get available cores according to the running OS
func OsCore() []proc {
	switch runtime.GOOS {
	case "linux":
		return []proc{{`NPROCS = $(shell grep -c 'processor' /proc/cpuinfo)`, `MAKEFLAGS += -j$(NPROCS)`}}
	case "darwin":
		return []proc{{`NPROCS = $(sysctl -n hw.ncpu')`, `MAKEFLAGS += -j$(NPROCS)`}}
	case "windows":
		return []proc{{`NPROCS = $(shell echo %NUMBER_OF_PROCESSORS%)`, `MAKEFLAGS += -j$(NPROCS)`}}
	default:
		return []proc{}
	}
}
