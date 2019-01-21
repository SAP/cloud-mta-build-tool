package proc

import "runtime"

// Proc - platform dependent commands and flags
type Proc struct {
	NPROCS    string
	MAKEFLAGS string
}

// OsCore - Get available cores according to the running OS
func OsCore() Proc {
	osProcMap := map[string]Proc{
		"linux":   {`NPROCS = $(shell grep -c 'processor' /proc/cpuinfo)`, `MAKEFLAGS += -j$(NPROCS)`},
		"darwin":  {`NPROCS = $(sysctl -n hw.ncpu)`, `MAKEFLAGS += -j$(NPROCS)`},
		"windows": {`NPROCS = $(shell echo %NUMBER_OF_PROCESSORS%)`, `MAKEFLAGS += -j$(NPROCS)`},
	}
	return osProcMap[runtime.GOOS]
}
