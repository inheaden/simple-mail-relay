PROJECT_NAME := "service"
PKG := "."
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all test clean grpc

all: build

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}

build: ## build the files
	@go build -o $(PROJECT_NAME) -i -v $(PKG)

build-static: ## Build the binary file
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o $(PROJECT_NAME) -installsuffix cgo -ldflags '-extldflags "-static"' -i -v $(PKG)

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'