.DEFAULT_GOAL := help
SHELL := /bin/bash

help:  # Print list of Makefile targets
	@for f in $(MAKEFILE_LIST); do grep -E ':  #' $${f} | grep -v 'LIST\|BEGIN' | \
	sort -u | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'; done
