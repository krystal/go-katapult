GOPATH ?= $(HOME)/go

BINDIR := bin
GOBIN_HELPER := github.com/myitcv/gobin
SOURCES := $(shell find . -name "*.go" -or -name "go.mod" -or -name "go.sum")

TOOLS += github.com/golangci/golangci-lint/cmd/golangci-lint@v1.31.0

export GOBIN := $(CURDIR)/$(BINDIR)
export PATH := $(GOBIN):$(PATH)
export GO111MODULE=on

ifdef VERBOSE
V = -v
endif

.PHONY: bootstrap
bootstrap: tools

.PHONY: test
test:
	go test $(V) -count=1 -race ./...

.PHONY: test-deps
test-deps:
	go test all

.PHONY: lint
lint:
	$(info Running Go linters)
	GOGC=off $(GOBIN)/golangci-lint $(V) run

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

.PHONY: clean
clean:
	rm -rf $(GOBIN)
	rm -f ./coverage.out ./go.mod.tidy-check ./go.sum.tidy-check

.PHONY: cov
cov: coverage.out

.PHONY: cov-html
cov-html: coverage.out
	go tool cover -html=./coverage.out

.PHONY: cov-func
cov-func: coverage.out
	go tool cover -func=./coverage.out

.PHONY: tools
tools: $(GOBIN_HELPER) $(TOOLS) $(TOOLS_INTERNAL)

.PHONY: $(TOOLS)
$(TOOLS): %:
	$(GOBIN)/gobin "$*"

.PHONY: $(GOBIN_HELPER)
$(GOBIN_HELPER): %:
	GO111MODULE=off go get -u $(GOBIN_HELPER)

.PHONY: deps
deps:
	$(info Downloading dependencies)
	go mod download

coverage.out: $(SOURCES)
	go test $(V) -covermode=count -coverprofile=./coverage.out ./...
