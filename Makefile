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
$(eval $(call tool,golangci-lint,github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51))
$(eval $(call tool,gomod,github.com/Helcaraxan/gomod@latest))

.PHONY: tools
tools: $(TOOLS)

#
# Generate
#

.PHONY: refresh-schema
refresh-schema:
	@curl https://my.katapult.io/core/v1/schema/openapi.json > katapult-openapi.json
	

.PHONY: generate
generate: 
	@docker run \
  --user "$$(id -u):$$(id -g)" \
  -v "$$(pwd):/local" \
  --rm \
  "$$(docker build -f container_images/oapi-codegen.Dockerfile --build-arg GO_VERSION="$$(grep goVersion generator-config.yml | cut -d':' -f2 | tr -d '[:space:]')" --build-arg GENERATOR_VERSION="$$(grep generatorVersion generator-config.yml | cut -d':' -f2 | tr -d '[:space:]')" -q .)" \
  -generate types,client \
  -package katapult \
  -templates /local/templates \
   /local/katapult-openapi.json > client.go


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
