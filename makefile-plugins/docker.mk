.DEFAULT_GOAL := help
SHELL := /bin/bash

.PHONY: docker-local
docker-local:  # Build a docker image for local development.
	docker build -t $(REPO_NAME):$(USER) .
	docker push $(REPO_NAME):$(USER)

.PHONY: docker-build
docker-build:  # Build a docker image without pushing it out.
	docker build -t $(DOCKER_USER)/$(REPO_NAME):$(DOCKER_TAG) .

.PHONY: docker-push
docker-push:  # Build a docker image without pushing it out.
	docker push $(DOCKER_USER)/$(REPO_NAME):$(DOCKER_TAG)

.PHONY: docker-all
docker-all: docker-build docker-push  # Build a docker image, tag, and push.
