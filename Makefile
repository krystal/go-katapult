GOMODNAME := $(shell grep 'module' go.mod | sed -e 's/^module //')
SOURCES := $(shell find . -name "*.go" -or -name "go.mod" -or -name "go.sum" \
	-or -name "Makefile")

# Verbose output
ifdef VERBOSE
V = -v
endif

# Global environment variables for all targets
SHELL ?= /bin/bash
SHELL := env \
	GO111MODULE=on \
	CGO_ENABLED=1 \
	$(SHELL)

#
# Defaults
#

# Default target
.DEFAULT_GOAL := test

#
# Development
#

BENCH ?= .
TESTARGS ?=
TESTTARGET ?= ./...

.PHONY: clean
clean:
	rm -rf ./bin/tools
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
lint:
	golangci-lint $(V) run --timeout=5m

.PHONY: format
format:
	git ls-files -z -- '*.go' ':(exclude)next/**/*.go' | xargs -0 goimports -w
	git ls-files -z -- '*.go' ':(exclude)next/**/*.go' | xargs -0 gofumpt -w

.PHONY: format-check
format-check:
	@unformatted="$$( \
		{ \
			git ls-files -z -- '*.go' ':(exclude)next/**/*.go' | \
				xargs -0 goimports -l; \
			git ls-files -z -- '*.go' ':(exclude)next/**/*.go' | \
				xargs -0 gofumpt -l; \
		} | sort -u \
	)"; \
	if [ -n "$$unformatted" ]; then \
		echo "Go files need formatting:" >&2; \
		echo "$$unformatted" >&2; \
		exit 1; \
	fi

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
	( diff -rN -x .git -x bin -x coverage.out \
		-x go.mod.tidy-check -x go.sum.tidy-check \
		"$(CURDIR)" "$(CHKDIR)" && rm -rf "$(CHKDIR)" ) || \
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
deps-analyze:
	gomod analyze

.PHONY: tidy
tidy:
	go mod tidy $(V)

.PHONY: verify
verify:
	go mod verify

.PHONY: check-tidy
check-tidy:
	go mod tidy -diff

#
# Documentation
#

# Serve docs
.PHONY: docs
docs:
	$(info serving docs on http://127.0.0.1:6060/pkg/$(GOMODNAME)/)
	@godoc -http=127.0.0.1:6060
