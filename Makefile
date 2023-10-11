# Makefile for running golangci-lint

GOLANGCI_VERSION := v1.51.2

.PHONY: lint

lint:
	@echo "Running golangci-lint..."
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:$(GOLANGCI_VERSION) golangci-lint run -c .golangci.yml

.PHONY: all

all: lint