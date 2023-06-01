# Make build for local usage
# The artifact is bin for each OS and copied to the go/bin dir
# Execute go generate to generate files to *.go to inclue as binary
# Execute go build
# Copy files to machine go/bin folder (temp target to avoid manual steps when developing locally)

all:format clean dir gen build-linux build-linux-arm build-darwin build-darwin-arm build-windows install-pkg install-micromatch-wrapper copy tests
.PHONY: build-darwin-arm build-darwin build-linux build-linux-arm build-windows tests

GOCMD=go
GOBUILD=$(GOCMD) build
GOLANGCI_VERSION = 1.21.0

# Binary names
BINARY_NAME=mbt
BUILD  = $(CURDIR)/release

# pkg
PKG_NAME=pkg

# micromatch wrapper
MICROMATCH_WRAPPER_DIR=$(CURDIR)/micromatch
MICROMATCH_WRAPPER_BINARY_NAME=micromatch-wrapper

ifeq ($(OS),Windows_NT)
MICROMATCH_WRAPPER_OS=win
else ifeq ($(shell uname -s), Linux)
MICROMATCH_WRAPPER_OS=linux
else ifeq ($(shell uname -s), Darwin)
MICROMATCH_WRAPPER_OS=macos
endif

ifeq ($(OS),Windows_NT)
	MICROMATCH_WRAPPER_SUFFIX = .exe
else
	MICROMATCH_WRAPPER_SUFFIX =
endif

format :
	go fmt ./...

tools:
	@echo "download golangci-lint"
	curl -sLO https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_VERSION}/golangci-lint-${GOLANGCI_VERSION}-linux-amd64.tar.gz
	tar -xzvf golangci-lint-${GOLANGCI_VERSION}-linux-amd64.tar.gz
	cp golangci-lint-${GOLANGCI_VERSION}-linux-amd64/golangci-lint $(GOPATH)/bin
	@echo "done"

lint:
	@echo "Start project linting"
	golangci-lint run --config .golangci.yml
	@echo "done linting"

# execute general tests
tests:
	 go test -v ./...
# check code coverage
cover:
	go test -v -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html
	open cover.html

clean:
	rm -rf $(BUILD)

dir:
	mkdir $(BUILD)

gen:
	go generate

# build for each platform
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_NAME)_linux -v

build-linux-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o release/$(BINARY_NAME)_linux_arm -v

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_NAME) -v

build-darwin-arm:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) -o release/$(BINARY_NAME)_darwin_arm -v

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_NAME)_windows -v

# use for local development - > copy the new bin to go/bin path to use new compiled version

# build and install micromatch wrapper
install-pkg:
	npm install -g $(PKG_NAME)
	echo "$(PKG_NAME) version:"
	$(PKG_NAME) --version

install-micromatch-wrapper:
	@cd $(MICROMATCH_WRAPPER_DIR) && npm install && cd -
	pkg $(MICROMATCH_WRAPPER_DIR) --out-path $(MICROMATCH_WRAPPER_DIR)
	cp $(MICROMATCH_WRAPPER_DIR)/$(MICROMATCH_WRAPPER_BINARY_NAME)-$(MICROMATCH_WRAPPER_OS)$(MICROMATCH_WRAPPER_SUFFIX) $(GOPATH)/bin/$(MICROMATCH_WRAPPER_BINARY_NAME)$(MICROMATCH_WRAPPER_SUFFIX)
	echo "$(MICROMATCH_WRAPPER_BINARY_NAME) version:"
	$(MICROMATCH_WRAPPER_BINARY_NAME) --version

copy:
ifeq ($(OS),Windows_NT)
	cp $(CURDIR)/release/$(BINARY_NAME)_windows $(GOPATH)/bin/$(BINARY_NAME).exe
else
	cp $(CURDIR)/release/$(BINARY_NAME) $(GOPATH)/bin/
	cp $(CURDIR)/release/$(BINARY_NAME) $~/usr/local/bin/
endif