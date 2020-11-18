SOURCES := $(shell find . -name "*.go" -or -name "go.mod" -or -name "go.sum" \
	-or -name "Makefile")

# Verbose output
ifdef VERBOSE
V = -v
endif

#
# Environment
#

BINDIR := bin
TOOLDIR := $(BINDIR)/tools

# Global environment variables for all targets
SHELL ?= /bin/bash
SHELL := env \
	GO111MODULE=on \
	GOBIN=$(CURDIR)/$(TOOLDIR) \
	CGO_ENABLED=1 \
	PATH='$(CURDIR)/$(BINDIR):$(CURDIR)/$(TOOLDIR):$(PATH)' \
	$(SHELL)

#
# Defaults
#

# Default target
.DEFAULT_GOAL := test

#
# Tools
#

TOOLS += $(TOOLDIR)/gobin
gobin: $(TOOLDIR)/gobin
$(TOOLDIR)/gobin:
	GO111MODULE=off go get -u github.com/myitcv/gobin

# external tool
define tool # 1: binary-name, 2: go-import-path
TOOLS += $(TOOLDIR)/$(1)

.PHONY: $(1)
$(1): $(TOOLDIR)/$(1)

$(TOOLDIR)/$(1): $(TOOLDIR)/gobin Makefile
	gobin $(V) "$(2)"
endef

$(eval $(call tool,gofumports,mvdan.cc/gofumpt/gofumports))
$(eval $(call tool,golangci-lint,github.com/golangci/golangci-lint/cmd/golangci-lint@v1.31))
$(eval $(call tool,golds,go101.org/golds@v0.1))

.PHONY: tools
tools: $(TOOLS)

#
# Development
#

TEST ?= $$(go list ./... | grep -v 'vendor')
BENCH ?= .

.PHONY: clean
clean:
	rm -f $(TOOLS)
	rm -f ./coverage.out ./go.mod.tidy-check ./go.sum.tidy-check

.PHONY: clean-golden
clean-golden:
	rm -f $(shell find * -path '*/testdata/*' -name "*.golden")

.PHONY: test
test:
	go test $(V) -count=1 -race $(TESTARGS) $(TEST)

.PHONY: test-update-golden
test-update-golden:
	@$(MAKE) test UPDATE_GOLDEN=1

.PHONY: regen-golden
regen-golden: clean-golden test-update-golden

.PHONY: test-deps
test-deps:
	go test all

.PHONY: lint
lint: golangci-lint
	GOGC=off golangci-lint $(V) run

.PHONY: format
format: gofumports
	gofumports -w .

.SILENT: bench
.PHONY: bench
bench:
	go test $(V) -count=1 -bench=$(BENCH) $(TESTARGS) $(TEST)

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

#
# Documentation
#

# Serve docs
.PHONY: docs
docs: golds
	$(info Serving go docs)
	@golds -emphasize-wdpkgs ./...
