LAMBOINJECTOR := lamboInjector
PKG_LIST := main.go
BUILD_DIR := "build"
#GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
COMMIT = $(shell git rev-list -1 HEAD)
# VERSION=$(shell cat ./version)

.PHONY: all dep build clean test coverage coverhtml lint

all: build

lint: ## Lint the files
	@$(HOME)/go/bin/golint -set_exit_status ${PKG_LIST}

dep: ## Get the dependencies
	go mod tidy

build: dep ## Build the binary file
	mkdir -p $(BUILD_DIR)
	go build -o build/$(LAMBOINJECTOR) $(PKG_LIST)

clean: ## Remove previous build
	@rm -fr $(BUILD_DIR)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'