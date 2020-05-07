# Copy this file to your project and customize it.

.DEFAULT_GOAL := help

ACTUAL_PWD := $(PWD)
PLUGINS_DIR := ./makefile-plugins
GO_SOURCE_PATHS := ./ cmd/... info/... pkg/... wrap/...

# These must be included first.
include $(PLUGINS_DIR)/common.mk

$(shell mkdir -p $(BUILD_DIR))
$(shell mkdir -p $(CACHE_DIR))

PROGRAM    := github-migration-cli
LICENSE    := Apache-License-2.0
PACKAGE    := github.com/parinithshekar/$(PROGRAM)
URL        := https://$(PACKAGE)
DOCKER_TAG := $(GIT_HASH)
TARGETS    := darwin/amd64 linux/amd64 windows/amd64
REPO_NAME  := github-migration-cli

GOFLAGS     := GOFLAGS="-mod=vendor"  # For vendored deps.
CGO_ENABLED := 0

UNAME_S?= $(shell uname -s)
ifeq ("${UNAME_S}", "Darwin")
    GOOS = darwin
else
    GOOS = linux
endif
GOARCH   := amd64

GOFILES := $(shell find . -type d \( -name .git -o -name vendor -o -name .submodules -o -name .cache \) -prune -o -type f -name "*.go" -print)

include $(PLUGINS_DIR)/docker.mk
include $(PLUGINS_DIR)/go.mk
# include $(PLUGINS_DIR)/swagger.mk

.PHONY: all
all:  # Run generally applicable Makefile targets (w/o lint).
all: vendor deps format test report build

.PHONY: all-w-lint
all-w-lint:  # Run generally applicable Makefile targets (w/lint).
all-w-lint: vendor deps format test report lint-go build

# TODO change makefile-plugins
.PHONY: build
build: build-linux # Build for target platforms.

.PHONY: clean
clean:  # Clean temporary files.
	@echo "==> Cleaning temporary files."
	rm -f ./*.html
	rm -f ./*.pdf
	rm -f ./*.pprof
	rm -f ./cp.out
	rm -f ./$(PROGRAM)
	rm -rf ./.cache/*
	rm -rf ./bin/*
	rm -rf ./build/*
	rm -rf ./dist/*

.PHONY: deps
deps:  # Install dependencies for development.
deps: deps-go # deps-swagger

# .PHONY: gen
# gen: models/action.go  # Regenerate generated artifacts.
# models/action.go: api/swagger/meraki-swagger.yaml
# 	./bin/swagger generate client -f ./api/swagger/meraki-swagger.yaml --name $(PROGRAM)

# .PHONY: validate
# validate:
# 	./bin/swagger validate ./api/swagger/meraki-swagger.yaml

# .PHONY: demo
# demo: ## run and record the demo-magic script
# 	which ttyrec || (sudo apt-get update && sudo apt-get install ttyrec)
# 	ttyrec -e './demo/demo.sh' ./demo/recording.ttyrec
# 	./bin/ttyrec2gif -in ./demo/recording.ttyrec -out demo/demo.gif -s 1.0 -col 120 -row 45
# 	rm -f ./demo/recording.ttyrec

# TODO This should be moved to makefile-plugins
.PHONY: dist
dist:  # Prepare artifacts under dist/
dist:
	@echo "==> Building all targets to dist/."
	@mkdir -p dist/
	@rm -rf dist/*
	@$(GOFLAGS) CGO_ENABLED=$(CGO_ENABLED) ./bin/gox \
		-parallel=3 \
		-output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" \
		-osarch='$(TARGETS)' \
		-gcflags '$(GCFLAGS)' \
		-ldflags '$(LDFLAGS)' \
		$(ACTUAL_PWD)
	@for D in dist/*; do \
		cp README.md $$D/; \
		T=$$(echo $$D | cut -f 2 -d /); \
		tar zcf dist/$$T.tar.gz -C $$D .; \
	done


include $(PLUGINS_DIR)/help.mk  # Must be included last.
