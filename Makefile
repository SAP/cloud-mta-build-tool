# Make build for local usage
# The artifact is bin for each OS and copied to the go/bin dir
# Execute go generate to generate files to *.go to inclue as binary
# Execute go build
# Copy files to machine go/bin folder (temp target to avoid manual steps when developing locally)

all:format clean dir gen build-linux build-darwin build-windows copy test
.PHONY: build-darwin build-linux build-windows


GOCMD=go
GOBUILD=$(GOCMD) build
GOLANGCI_VERSION = 1.21.0

# Binary names
BINARY_NAME=mbt
BUILD  = $(CURDIR)/release

format :
	go fmt ./...

tools_lcl:
	@echo "download golangci-lint"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest
#	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(go env GOPATH)/bin latest
#	golangci-lint --version

lint:
	@echo "Start project linting"
	golangci-lint run --config .golangci.yml
	@echo "done"

# execute general tests
test:
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

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_NAME) -v

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
