
# Create folder for build artifacts
DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BASE := $(notdir $(patsubst %/,%,$(DIR)))
MTR_DIR := $(DIR)$(BASE)
# Determine OS cores
NPROCS = $(sysctl -n hw.ncpu')
MAKEFLAGS += -j$(NPROCS)
# List of all the recipes to be executed during the build process
.PHONY: meta mtar cleanup
# Default target compile all
all: tdir meta mtar cleanup
# Create folder for build artifacts
tdir:
	@mkdir -p "$(MTR_DIR)"
# Set the current project repository path for general mta process
  PROJ_DIR := $(CURDIR)

modules := $(shell mbt provide modules)
modules := $(subst ],,$(subst [,,$(modules)))

##Execute module build
execute:
	@for module in $(modules) ; do \
	 mbt build module $$module ; \
	done
