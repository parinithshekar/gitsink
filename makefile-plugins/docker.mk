.DEFAULT_GOAL := help
SHELL := /bin/bash

TAG_NAME := $(DOCKER_TAG)
BASE_IMAGE := $(DOCKER_USER)/$(REPO_NAME)
NEW_IMAGE_TAG := $(BASE_IMAGE):$(TAG_NAME)

.PHONY: docker-login
docker-login: # Login to dockerhub account
	echo "$${DOCKER_PASS}" | docker login --username "$${DOCKER_USER}" --password-stdin

.PHONY: docker-local
docker-local:  # Build a docker image for local development.
	docker build -t $(BASE_IMAGE):latest .

.PHONY: docker-build
docker-build:  # Build a docker image without pushing it out.
	docker build -t $(BASE_IMAGE):latest .
	docker tag $(BASE_IMAGE):latest $(NEW_IMAGE_TAG)

.PHONY: docker-push
docker-push:  # Push docker image to dockerhub
	docker push $(BASE_IMAGE):latest
	docker push $(NEW_IMAGE_TAG)

.PHONY: docker-all
docker-all: docker-build docker-push  # Build a docker image, tag, and push.
