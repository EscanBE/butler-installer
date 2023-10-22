GIT_TAG := $(shell echo $(shell git describe --tags || git branch --show-current) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')
BUILD_DATE	:= $(shell date '+%Y-%m-%d')

###############################################################################
###                                Build flags                              ###
###############################################################################

LD_FLAGS = -X github.com/EscanBE/go-app-name/constants.VERSION=$(GIT_TAG) \
            -X github.com/EscanBE/go-app-name/constants.COMMIT_HASH=$(COMMIT) \
            -X github.com/EscanBE/go-app-name/constants.BUILD_DATE=$(BUILD_DATE)

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

###############################################################################
###                                  Test                                   ###
###############################################################################

test: go.sum
	@echo "testing"
	@go test -v ./... -race -coverprofile=coverage.txt -covermode=atomic
.PHONY: test

###############################################################################
###                                  Build                                  ###
###############################################################################

build: go.sum
ifeq ($(OS),Windows_NT)
	@echo "building goappnamed binary..."
	@echo "Flags $(BUILD_FLAGS)"
	@go build -mod=readonly $(BUILD_FLAGS) -o build/goappnamed.exe ./cmd/goappnamed
	@echo "Builded successfully"
else
	@echo "building goappnamed binary..."
	@echo "Flags $(BUILD_FLAGS)"
	@go build -mod=readonly $(BUILD_FLAGS) -o build/goappnamed ./cmd/goappnamed
	@echo "Builded successfully"
endif
.PHONY: build

###############################################################################
###                                 Install                                 ###
###############################################################################

install: go.sum
	@echo "installing goappnamed binary..."
	@echo "Flags $(BUILD_FLAGS)"
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/goappnamed
	@echo "Installed successfully"
.PHONY: install