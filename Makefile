# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=chromium-downloader
BINARY_WIN=$(BINARY_NAME)-windows-amd64
BINARY_UNIX=$(BINARY_NAME)-linux-amd64
BINARY_MAC=$(BINARY_NAME)-darwin-amd64

build:
	$(GOBUILD) -ldflags="-s -w" -o builds/$(BINARY_UNIX) -v

run:
	$(GOBUILD) -ldflags="-s -w" -o builds/$(BINARY_UNIX) -v ./...
	./builds/$(BINARY_UNIX)

# Cross compilation

build-linux:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o builds/$(BINARY_UNIX) -v

build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o builds/$(BINARY_UNIX) -v

build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o builds/$(BINARY_UNIX) -v


clean: 
	$(GOCLEAN)
	rm -f builds/$(BINARY_WIN)
	rm -f builds/$(BINARY_UNIX)
	rm -f builds/$(BINARY_MAC)

OUT = builds

.PHONY: all
all: build

$(shell mkdir -p $(OUT))