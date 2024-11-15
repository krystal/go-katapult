GOMODNAME := $(shell grep 'module' go.mod | sed -e 's/^module //')
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

# external tool
define tool # 1: binary-name, 2: go-import-path
TOOLS += $(TOOLDIR)/$(1)

$(TOOLDIR)/$(1): Makefile
	GOBIN="$(CURDIR)/$(TOOLDIR)" go install "$(2)"
endef

$(eval $(call tool,godoc,golang.org/x/tools/cmd/godoc@latest))
$(eval $(call tool,gofumpt,mvdan.cc/gofumpt@latest))
$(eval $(call tool,goimports,golang.org/x/tools/cmd/goimports@latest))
$(eval $(call tool,golangci-lint,github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62))
$(eval $(call tool,gomod,github.com/Helcaraxan/gomod@latest))

.PHONY: tools
tools: $(TOOLS)

#
# Development
#

BENCH ?= .
TESTARGS ?=
TESTTARGET ?= ./...

.PHONY: clean
clean:
	rm -f $(BINS) $(TOOLS)
	rm -f ./coverage.out ./go.mod.tidy-check ./go.sum.tidy-check

.PHONY: clean-golden
clean-golden:
	rm -f $(shell find * -path '*/testdata/*' -name "*.golden" \
		-exec echo "'{}'" \;)

.PHONY: test
test:
	go test $(V) -count=1 -race $(TESTARGS) $(TESTTARGET)

.PHONY: test-deps
test-deps:
	@$(MAKE) test TESTTARGET=all

.PHONY: lint
lint: $(TOOLDIR)/golangci-lint
	golangci-lint $(V) run --timeout=5m

.PHONY: format
format: $(TOOLDIR)/goimports $(TOOLDIR)/gofumpt
	goimports -w . && gofumpt -w .

.SILENT: bench
.PHONY: bench
bench:
	@$(MAKE) test TESTARGS="-bench=$(BENCH) -benchmem"

.PHONY: update-golden
update-golden:
	@$(MAKE) test GOLDEN_UPDATE=1

.PHONY: regen-golden
regen-golden: clean-golden update-golden

#
# Code Generation
#

.PHONY: generate
generate: schemas
	go generate ./...

.PHONY: check-generate
check-generate:
	$(eval CHKDIR := $(shell mktemp -d))
	cp -a . "$(CHKDIR)"
	make -C "$(CHKDIR)/" generate
	( diff -rN "$(CURDIR)" "$(CHKDIR)" && rm -rf "$(CHKDIR)" ) || \
		( rm -rf "$(CHKDIR)" && exit 1 )

.PHONY: generate-next
generate-next:
	( cd ./next && ./generate.sh )

#
# Katapult API Schemas
#

.PHONY: schemas
schemas:
	go generate ./schemas

.PHONY: update-schemas
update-schemas:
	SCHEMA_FORCE_UPDATE=1 go generate ./schemas

.PHONY: check-schemas
check-schemas:
	$(eval CHKDIR := $(shell mktemp -d))
	cp -a . "$(CHKDIR)"
	make -C "$(CHKDIR)/" update-schemas
	( diff -rN "$(CURDIR)/schemas" "$(CHKDIR)/schemas" && rm -rf "$(CHKDIR)" ) \
		|| ( rm -rf "$(CHKDIR)" && exit 1 )

.PHONY: retrieve-openapi-schemas
retrieve-openapi-schemas:
	wget -O next/katapult-core-openapi.json https://api.katapult.io/core/v1/schema/openapi.json
	wget -O next/katapult-public-openapi.json https://api.katapult.io/public/v1/schema/openapi.json

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
	@$(MAKE) test TESTARGS="-covermode=atomic -coverprofile=./coverage.out"

#
# Dependencies
#

.PHONY: deps
deps:
	$(info Downloading dependencies)
	go mod download

.PHONY: deps-update
deps-update:
	go get -u -t ./...

.PHONY: deps-analyze
deps-analyze: $(TOOLDIR)/gomod
	gomod analyze

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
	( \
		diff go.mod go.mod.tidy-check && \
		diff go.sum go.sum.tidy-check && \
		rm -f go.mod go.sum && \
		mv go.mod.tidy-check go.mod && \
		mv go.sum.tidy-check go.sum \
	) || ( \
		rm -f go.mod go.sum && \
		mv go.mod.tidy-check go.mod && \
		mv go.sum.tidy-check go.sum; \
		exit 1 \
	)

#
# Documentation
#

# Serve docs
.PHONY: docs
docs: $(TOOLDIR)/godoc
	$(info serviing docs on http://127.0.0.1:6060/pkg/$(GOMODNAME)/)
	@godoc -http=127.0.0.1:6060
