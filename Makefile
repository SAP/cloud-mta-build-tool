# Make for build locally
# Create build with OS artifact's which need to put under the bin file as executable bin

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test


# Binary names
BINARY_NAME=mit
BUILD  = $(CURDIR)/release


all:clean dir build-linux build-darwin build-windows copy
.PHONY: build-darwin build-linux build-windows

clean:
	rm -rf $(BUILD)

dir:
	mkdir $(BUILD)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_NAME)_linux -v

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_NAME) -v

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_NAME)_windows -v

# Use for local - > copy the bin to go/bin
copy:
	cp $(CURDIR)/release/$(BINARY_NAME) $(GOPATH)/bin/





