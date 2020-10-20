GOPATH ?= $(HOME)/go
SOURCES := $(shell find . -name "*.go" -or -name "go.mod" -or -name "go.sum")

TOOLSDIR := tools
export GOBIN := $(CURDIR)/$(TOOLSDIR)
export PATH := $(GOBIN):$(PATH)
export GO111MODULE=on

GOBIN_SRC := github.com/myitcv/gobin
GOBIN_BIN := $(GOBIN)/gobin

ifdef VERBOSE
V = -v
endif

# Default target
.DEFAULT_GOAL := test

#
# Tools
#

define tool # 1: binary-name, 2: go-src-path
TOOLS += $(GOBIN)/$(1)

$(GOBIN)/$(1): $(GOBIN_BIN)
	$(GOBIN_BIN) "$(2)"
endef

$(GOBIN_BIN): %:
	GO111MODULE=off go get -u $(GOBIN_SRC)

$(eval $(call tool,golangci-lint,github.com/golangci/golangci-lint/cmd/golangci-lint@v1.31.0))

.PHONY: tools
tools: $(GOBIN_BIN) $(TOOLS)

#
# Development
#

.PHONY: clean
clean:
	rm -rf $(GOBIN)
	rm -f ./coverage.out ./go.mod.tidy-check ./go.sum.tidy-check

.PHONY: test
test:
	go test $(V) -count=1 -race ./...

.PHONY: test-deps
test-deps:
	go test all

.PHONY: lint
lint: $(TOOLS)
	$(info Running Go linters)
	GOGC=off $(GOBIN)/golangci-lint $(V) run

#
# Coverage
#

.PHONY: cov
cov: coverage.out

.PHONY: cov-html
cov-html: coverage.out
	go tool cover -html=./coverage.out

.PHONY: cov-func
cov-func: coverage.out
	go tool cover -func=./coverage.out

coverage.out: $(SOURCES)
	go test $(V) -covermode=count -coverprofile=./coverage.out ./...

#
# Dependencies
#

.PHONY: deps
deps:
	$(info Downloading dependencies)
	go mod download

.PHONY: tidy
tidy:
	go mod tidy $(V)

.PHONY: verify
verify:
	go mod verify

.SILENT: check-tidy
.PHONY: check-tidy
check-tidy:
	cp go.mod go.mod.tidy-check
	cp go.sum go.sum.tidy-check
	go mod tidy
	-diff go.mod go.mod.tidy-check
	-diff go.sum go.sum.tidy-check
	-rm -f go.mod go.sum
	-mv go.mod.tidy-check go.mod
	-mv go.sum.tidy-check go.sum
