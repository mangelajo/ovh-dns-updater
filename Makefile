# Makefile for OVH DNS Updater

# Variables
APP_NAME := ovh-dns-updater
CONTAINER_REPO := # Add your container repository here
VERSION := $(shell git describe --tags --always --dirty)
HELM_CHART_PATH := charts/ovh-dns-updater
HELM_REPO := # Add your Helm repository here

# Container engine detection (docker or podman)
CONTAINER_ENGINE ?= $(shell command -v podman > /dev/null 2> /dev/null && echo podman || echo docker)

# Go related variables
GOFILES := $(wildcard *.go)
GOBIN := $(GOPATH)/bin

# Build the Go binary
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) -v

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME)
	rm -f coverage.out coverage.html

# Build container image
.PHONY: container-build
container-build:
	@echo "Building container image with $(CONTAINER_ENGINE)..."
	$(CONTAINER_ENGINE) build -t $(APP_NAME):$(VERSION) .
	$(CONTAINER_ENGINE) tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

# Push container image
.PHONY: container-push
container-push: container-build
	@echo "Pushing container image with $(CONTAINER_ENGINE)..."
	$(CONTAINER_ENGINE) tag $(APP_NAME):$(VERSION) $(CONTAINER_REPO)/$(APP_NAME):$(VERSION)
	$(CONTAINER_ENGINE) tag $(APP_NAME):$(VERSION) $(CONTAINER_REPO)/$(APP_NAME):latest
	$(CONTAINER_ENGINE) push $(CONTAINER_REPO)/$(APP_NAME):$(VERSION)
	$(CONTAINER_ENGINE) push $(CONTAINER_REPO)/$(APP_NAME):latest

# Lint Helm chart
.PHONY: helm-lint
helm-lint:
	@echo "Linting Helm chart..."
	helm lint $(HELM_CHART_PATH)

# Package Helm chart
.PHONY: helm-package
helm-package: helm-lint
	@echo "Packaging Helm chart..."
	helm package $(HELM_CHART_PATH) --version $(VERSION) --app-version $(VERSION)

# Push Helm chart to repository
.PHONY: helm-push
helm-push: helm-package
	@echo "Pushing Helm chart to repository..."
	helm push $(APP_NAME)-$(VERSION).tgz $(HELM_REPO)

# Install Helm chart locally (for testing)
.PHONY: helm-install
helm-install: helm-package
	@echo "Installing Helm chart locally..."
	helm install $(APP_NAME) ./$(APP_NAME)-$(VERSION).tgz

# Run the application locally
.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	./$(APP_NAME)

# All-in-one target for CI/CD
.PHONY: ci
ci: test container-build helm-package

# Release target - builds, tests, packages and pushes everything
.PHONY: release
release: test container-push helm-push
	@echo "Release $(VERSION) completed!"

# Default target
.PHONY: all
all: build test

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the Go binary"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Clean build artifacts"
	@echo "  container-build - Build container image with $(CONTAINER_ENGINE)"
	@echo "  container-push  - Build and push container image with $(CONTAINER_ENGINE)"
	@echo "  helm-lint      - Lint Helm chart"
	@echo "  helm-package   - Package Helm chart"
	@echo "  helm-push      - Package and push Helm chart"
	@echo "  helm-install   - Install Helm chart locally"
	@echo "  run            - Run the application locally"
	@echo "  ci             - Run CI tasks (test, docker-build, helm-package)"
	@echo "  release        - Release a new version (test, docker-push, helm-push)"
	@echo "  all            - Build and test"
	@echo "  help           - Show this help message"
