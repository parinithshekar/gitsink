.DEFAULT_GOAL := help
SHELL := /bin/bash

# helpers to find the current path if THIS makefile which may be imported
__go_mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
__go_mkfile_dir := $(dir $(__go_mkfile_path))

.PHONY: deps-swagger
deps-swagger:  # Install dependencies for Swagger development.
	@echo "==> Installing dependencies for Swagger development."
	$(__go_mkfile_dir)/swagger.sh install-go-swagger ./bin/swagger

gen-swagger:
	./bin/swagger validate api/swagger/meraki-swagger.yaml
