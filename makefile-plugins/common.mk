.DEFAULT_GOAL := help
SHELL := /bin/bash

DATE := $(shell date +%FT%T%z)
USER := $(shell whoami)
GIT_HASH := $(shell git --no-pager describe --tags --always)
BRANCH := $(shell git branch | grep '*' | cut -d ' ' -f2)

BUILD_DIR := $(ACTUAL_PWD)/build
CACHE_DIR := $(ACTUAL_PWD)/.cache
