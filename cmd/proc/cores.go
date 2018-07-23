package proc

import "runtime"

type Proc struct {
	NPROCS    string
	MAKEFLAGS string
}

// OsCore - Get the build operation's
func OsCore() []Proc {
	switch runtime.GOOS {
	case "linux":
		return []Proc{{`NPROCS = $(shell grep -c 'processor' /proc/cpuinfo)`, `MAKEFLAGS += -j$(NPROCS)`}}
	case "darwin":
		return []Proc{{`NPROCS = $(sysctl -n hw.ncpu')`, `MAKEFLAGS += -j$(NPROCS)`}}
	case "windows":
		return []Proc{}
	default:
		return []Proc{}
	}
}
