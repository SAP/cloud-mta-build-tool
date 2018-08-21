# Make build for local usage
# The artifact is bin for each OS and copied to the go/bin dir
# Execute go generate to generate files to *.go to inclue as binary
# Execute go build
# Copy files to machine go/bin folder (temp target to avoid manual steps when developing locally)

GOCMD=go
GOBUILD=$(GOCMD) build


# Binary names
BINARY_NAME=mbt
BUILD  = $(CURDIR)/release


all:clean dir gen build-linux build-darwin build-windows copy
.PHONY: build-darwin build-linux build-windows

clean:
	rm -rf $(BUILD)

dir:
	mkdir $(BUILD)

lint :
	@echo "style linting files:" # this list will grow as we cleanup all the code
	@bash -c "go run bin/style/main.go api util/{configv3,manifest,randomword,sorting,ui}"
	@echo "No lint errors!"
	@echo

gen:
	go generate

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_NAME)_linux -v

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_NAME) -v

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_NAME)_windows -v

# Use for local - > copy the bin to go/bin
copy:
	cp $(CURDIR)/release/$(BINARY_NAME) $(GOPATH)/bin/







