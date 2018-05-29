GOFILES := $(shell find . -type f -name '*.go')
TARGET := secret-sender
VERSION=`git rev-parse HEAD`
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION}"
BUILD_NAME := github.com/Shopify/secret-sender

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(dir $(mkfile_path))

all: default
default: secret-sender

$(TARGET): $(GOFILES)
	@echo "   \x1b[1;34mgo build\x1b[0m  $@"
	@env GOBIN=$(current_dir) go install $(LDFLAGS) $(BUILD_NAME)

clean:
	@echo "         \x1b[1;31mrm\x1b[0m  $(TARGET)"
	@rm -f $(TARGET)

.PHONY: all default clean
