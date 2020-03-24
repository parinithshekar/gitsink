.DEFAULT_GOAL := help
SHELL := /bin/bash

INFO_PACKAGE := $(PACKAGE)/info

LDFLAGS := -X "$(INFO_PACKAGE).Program=$(PROGRAM)"
LDFLAGS += -X "$(INFO_PACKAGE).License=$(LICENSE)"
LDFLAGS += -X "$(INFO_PACKAGE).URL=$(URL)"
LDFLAGS += -X "$(INFO_PACKAGE).BuildUser=$(USER)"
LDFLAGS += -X "$(INFO_PACKAGE).BuildDate=$(DATE)"
LDFLAGS += -X "$(INFO_PACKAGE).Version=$(GIT_HASH)"
LDFLAGS += -X "$(INFO_PACKAGE).Revision=$(GIT_HASH)"
LDFLAGS += -X "$(INFO_PACKAGE).Branch=$(BRANCH)"

#######################################
# Disable Golang Debugger (strip symbols from binary)
#   Activate the following lines instead
# LDFLAGS += -s
# LDFLAGS += -linkmode external -extldflags -static -s -w
#######################################
# Enable Golang Debugger (keep symbols)
LDFLAGS += -extldflags -static
# These flags disable compiler optimizations so that debuggers work
GCFLAGS := -N -l
#######################################

# helpers to find the current path if THIS makefile which may be imported
__go_mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
__go_mkfile_dir := $(dir $(__go_mkfile_path))

.PHONY: deps-go
deps-go:  # Install dependencies for Go development.
	@echo "==> Installing dependencies for Go development."
	$(__go_mkfile_dir)/go.sh install-dlv ./bin/dlv
	$(__go_mkfile_dir)/go.sh install-golangci-lint ./bin/golangci-lint
	$(__go_mkfile_dir)/go.sh install-gox ./bin/gox
	$(__go_mkfile_dir)/go.sh install-ttyrec2gif ./bin/ttyrec2gif

.PHONY: lint-go
lint-go:  # go lint check
	@golint $(GO_SOURCE_PATHS)

.PHONY: build-linux
build-linux:  # Build for Linux.
	@echo "==> Building for Linux."
	$(GOFLAGS) CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
	  go build -a -gcflags '$(GCFLAGS)' -ldflags '$(LDFLAGS)' $(ACTUAL_PWD)

.PHONY: format
format:  # Auto-format all Go files not in vendor dir.
	@echo "==> Auto-formatting all Go files not in vendor dir."
	gofmt -s -w $(GOFILES)

.PHONY: lint
lint:  # Lint all Go source code using GolangCI-Lint.
	@echo "==> Linting all Go source code using GolangCI-Lint."
	./bin/golangci-lint run \
	  --modules-download-mode vendor \
	  --skip-dirs '(vendor|.submodules|.cache|.git)' \
	  --disable errcheck \
	  --enable gofmt \
	  --enable goimports \
	  --timeout 600s \
	  -v

.PHONY: report
report:  # Generate all reports.
	@echo "==> Generating profiler reports."
	for mode in cpu mem mutex block; do \
	  if [ -e $$mode.pprof ]; then \
	    go tool pprof --pdf $(PROGRAM) \$$mode.pprof > $$mode.pprof.pdf; \
	  fi; done
	@if [ -d ~/x/tmp ] && compgen -G "*.pdf"; then cp -v *.pdf ~/x/tmp; fi
	@echo "==> Generating coverage reports."
	$(GOFLAGS) go tool cover -html=cp.out -o=coverage.html
	@if [ -d ~/x/tmp ] && compgen -G "*.html"; then cp -v *.html ~/x/tmp; fi

.PHONY: test
test:  # Run all tests.
	@echo "==> Running all tests."
	$(GOFLAGS) go test ./... -coverprofile=cp.out
	$(GOFLAGS) go tool cover -func=cp.out 2>&1 | tee .cache/test_coverage.txt
	cat .cache/test_coverage.txt | grep total | grep -Eo '[0-9]+\.[0-9]+' > .cache/test_coverage_total.txt

.PHONY: vendor
vendor: __update-package-deps  # Vendor all Go dependencies.

.PHONY: __update-package-deps
__update-package-deps:
	@if ! $(__go_mkfile_dir)/go.sh package-deps-unchanged? \
	        $(shell find . -type d \( -name .git -o -name vendor -o -name .submodules -o -name .cache \) -prune \
	        -o -type f -name "*.go" -print | sed 's@^./@@g' | cut -d'/' -f1 | sort | uniq) ; \
	then \
	  echo "==> Vendoring all Go dependencies."; \
	  go mod tidy && \
	  go mod vendor && \
	  go mod verify; \
	fi
