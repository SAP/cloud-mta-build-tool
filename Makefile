# Make build for local usage
# The artifact is bin for each OS and copied to the go/bin dir
# Execute go generate to generate files to *.go to inclue as binary
# Execute go build
# Copy files to machine go/bin folder (temp target to avoid manual steps when developing locally)

all:format clean dir gen build-linux build-linux-arm build-darwin build-darwin-arm build-windows copy install-cyclonedx tests
.PHONY: build-darwin-arm build-darwin build-linux build-linux-arm build-windows tests

GOCMD=go
GOBUILD=$(GOCMD) build
GOLANGCI_VERSION = 1.21.0

# Binary names
BINARY_NAME=mbt
BUILD  = $(CURDIR)/release

# cyclonedx-cli
CYCLONEDX_CLI_BINARY = cyclonedx
CYCLONEDX_CLI_VERSION = 0.24.2

# cyclonedx-gomod
CYCLONEDX_GOMOD_BINARY = cyclonedx-gomod
CYCLONEDX_GOMOD_VERSION = latest

# cyclonedx-bom
CYCLONEDX_BOM_PACKAGE = cyclonedx-bom
CYCLONEDX_BOM_VERSION = 0.0.9
CYCLONEDX_BOM_BINARY = cyclonedx-bom


ifeq ($(OS),Windows_NT)
CYCLONEDX_OS=win
else ifeq ($(shell uname -s), Linux)
CYCLONEDX_OS=linux
else ifeq ($(shell uname -s), Darwin)
CYCLONEDX_OS=osx
endif

ifeq ($(shell uname -m),x86_64)
	CYCLONEDX_ARCH=x64
else ifeq ($(shell uname -m),arm64)
	CYCLONEDX_ARCH=arm64
else ifeq ($(shell uname -m),i386)
	CYCLONEDX_ARCH=x86
else
	CYCLONEDX_ARCH=x64
endif

ifeq ($(OS),Windows_NT)
	CYCLONEDX_BINARY_SUFFIX = .exe
else
	CYCLONEDX_BINARY_SUFFIX =
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
	 go test -v -count=1 -timeout 30m ./...
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
copy:
ifeq ($(OS),Windows_NT)
	cp $(CURDIR)/release/$(BINARY_NAME)_windows $(GOPATH)/bin/$(BINARY_NAME).exe
else
	cp $(CURDIR)/release/$(BINARY_NAME) $(GOPATH)/bin/
	cp $(CURDIR)/release/$(BINARY_NAME) $~/usr/local/bin/
endif

# use for local development - > install cyclonedx-gomod, cyclonedx-cli and cyclonedx-bom
install-cyclonedx:
# install cyclonedx-gomod
	go install github.com/CycloneDX/cyclonedx-gomod/cmd/${CYCLONEDX_GOMOD_BINARY}@${CYCLONEDX_GOMOD_VERSION}
	echo "${CYCLONEDX_GOMOD_BINARY} version"
	${CYCLONEDX_GOMOD_BINARY} version
# install cyclonedx-cli	
	curl -fsSLO --compressed "https://github.com/CycloneDX/cyclonedx-cli/releases/download/v${CYCLONEDX_CLI_VERSION}/${CYCLONEDX_CLI_BINARY}-${CYCLONEDX_OS}-${CYCLONEDX_ARCH}${CYCLONEDX_BINARY_SUFFIX}"
	mv ${CYCLONEDX_CLI_BINARY}-${CYCLONEDX_OS}-${CYCLONEDX_ARCH}${CYCLONEDX_BINARY_SUFFIX} $(GOPATH)/bin/${CYCLONEDX_CLI_BINARY}${CYCLONEDX_BINARY_SUFFIX}
	echo "${CYCLONEDX_CLI_BINARY} version:"
	${CYCLONEDX_CLI_BINARY} --version
# install cyclonedx-bom
	npm install -g ${CYCLONEDX_BOM_PACKAGE}@${CYCLONEDX_BOM_VERSION}
	echo "${CYCLONEDX_BOM_BINARY} -h"
	npx ${CYCLONEDX_BOM_BINARY} -h