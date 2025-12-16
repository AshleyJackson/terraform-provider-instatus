# Terraform Provider for Instatus - Makefile

# Detect OS and architecture
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_S),Darwin)
    OS := darwin
endif
ifeq ($(UNAME_S),Linux)
    OS := linux
endif

ifeq ($(UNAME_M),arm64)
    ARCH := arm64
endif
ifeq ($(UNAME_M),x86_64)
    ARCH := amd64
endif

BINARY_NAME := terraform-provider-instatus
INSTALL_PATH := ~/.terraform.d/plugins/registry.terraform.io/ashleyjackson/instatus/0.1.0/$(OS)_$(ARCH)

.PHONY: help build install test clean fmt vet

help: ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the provider binary
	@echo "Building $(BINARY_NAME) for $(OS)_$(ARCH)..."
	go build -o $(BINARY_NAME)
	@echo "Build complete!"

install: build ## Build and install the provider locally
	@echo "Installing provider to $(INSTALL_PATH)..."
	mkdir -p $(INSTALL_PATH)
	cp $(BINARY_NAME) $(INSTALL_PATH)/
	@echo "Provider installed successfully!"
	@echo "Run 'terraform init' in your project to use it."

test: ## Run tests
	go test -v ./...

fmt: ## Format Go code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

tidy: ## Tidy Go modules
	go mod tidy

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -rf test/.terraform test/.terraform.lock.hcl test/terraform.tfstate*
	@echo "Clean complete!"

uninstall: ## Uninstall the provider
	@echo "Uninstalling provider from $(INSTALL_PATH)..."
	rm -rf $(INSTALL_PATH)
	@echo "Provider uninstalled!"

test-setup: install ## Install provider and initialize test directory
	@echo "Setting up test environment..."
	cd test && rm -rf .terraform .terraform.lock.hcl && terraform init
	@echo "Test environment ready!"

test-apply: test-setup ## Run terraform apply in test directory
	@echo "Running terraform apply..."
	cd test && terraform apply -auto-approve

test-destroy: ## Run terraform destroy in test directory
	@echo "Running terraform destroy..."
	cd test && terraform destroy -auto-approve

all: fmt vet build install ## Format, vet, build, and install

.DEFAULT_GOAL := help
