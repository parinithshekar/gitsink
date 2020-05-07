.DEFAULT_GOAL := help
SHELL := /bin/bash

.PHONY: docker-login
docker-login: # Login to dockerhub account
	echo "$${DOCKER_PASS}" | docker login --username "$${DOCKER_USER}" --password-stdin

.PHONY: docker-local
docker-local:  # Build a docker image for local development.
	docker build -t $(REPO_NAME):$(USER) .
	docker push $(REPO_NAME):$(USER)

.PHONY: docker-build
docker-build:  # Build a docker image without pushing it out.
	docker build -t $(DOCKER_USER)/$(REPO_NAME):$(DOCKER_TAG) .

.PHONY: docker-push
docker-push:  # Push docker image to dockerhub
	docker push $(DOCKER_USER)/$(REPO_NAME):$(DOCKER_TAG)

.PHONY: docker-all
docker-all: docker-build docker-push  # Build a docker image, tag, and push.
