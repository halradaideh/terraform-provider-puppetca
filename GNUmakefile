TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=puppetca
PLUGIN_NAME=terraform-provider-puppetca
PLUGIN_VERSION=1.0.0

default: build

build: fmtcheck
	go install

# Build binary locally for installation
build-local: fmtcheck
	go build -o $(PLUGIN_NAME)

# Install to the legacy plugin path for local testing/development
local-install: build-local
	@echo "Installing provider to ~/.terraform.d/plugins/puppetca_v$(PLUGIN_VERSION)/$(PLUGIN_NAME)"
	@mkdir -p ~/.terraform.d/plugins/puppetca_v$(PLUGIN_VERSION)
	@cp $(PLUGIN_NAME) ~/.terraform.d/plugins/puppetca_v$(PLUGIN_VERSION)/$(PLUGIN_NAME)
	@echo "Local install complete! Legacy plugin path: ~/.terraform.d/plugins/puppetca_v$(PLUGIN_VERSION)/$(PLUGIN_NAME)"

# Install to the modern plugin path for Terraform 0.13+
local-install-modern: build-local
	$(eval OS_ARCH := $(shell go env GOOS)_$(shell go env GOARCH))
	@echo "Installing provider to ~/.terraform.d/plugins/registry.terraform.io/camptocamp/puppetca/$(PLUGIN_VERSION)/$(OS_ARCH)/$(PLUGIN_NAME)_v$(PLUGIN_VERSION)"
	@mkdir -p ~/.terraform.d/plugins/registry.terraform.io/camptocamp/puppetca/$(PLUGIN_VERSION)/$(OS_ARCH)
	@cp $(PLUGIN_NAME) ~/.terraform.d/plugins/registry.terraform.io/camptocamp/puppetca/$(PLUGIN_VERSION)/$(OS_ARCH)/$(PLUGIN_NAME)_v$(PLUGIN_VERSION)
	@echo "Modern install complete! Plugin path: ~/.terraform.d/plugins/registry.terraform.io/camptocamp/puppetca/$(PLUGIN_VERSION)/$(OS_ARCH)/$(PLUGIN_NAME)_v$(PLUGIN_VERSION)"

# Install to the local development path (matches local/puppetca/puppetca source)
local-install-dev: build-local
	$(eval OS_ARCH := $(shell go env GOOS)_$(shell go env GOARCH))
	@echo "Installing provider to ~/.terraform.d/plugins/local/puppetca/puppetca/$(PLUGIN_VERSION)/$(OS_ARCH)/$(PLUGIN_NAME)"
	@mkdir -p ~/.terraform.d/plugins/local/puppetca/puppetca/$(PLUGIN_VERSION)/$(OS_ARCH)
	@cp $(PLUGIN_NAME) ~/.terraform.d/plugins/local/puppetca/puppetca/$(PLUGIN_VERSION)/$(OS_ARCH)/$(PLUGIN_NAME)
	@echo "Development install complete! Plugin path: ~/.terraform.d/plugins/local/puppetca/puppetca/$(PLUGIN_VERSION)/$(OS_ARCH)/$(PLUGIN_NAME)"

# Install to all plugin paths (legacy, modern, and development)
install-local: local-install local-install-modern local-install-dev
	@echo "Provider installed to legacy, modern, and development plugin paths"

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v -parallel 20 $(TESTARGS) -timeout 120m

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w $(GOFMT_FILES)

# Currently required by tf-deploy compile
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run ./...

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

vendor:
	go mod tidy
	go mod vendor

vet:
	go vet ./...

# CI/CD targets
ci-test: fmtcheck lint vet test
	@echo "==> All CI tests passed!"

ci-build: fmtcheck
	@echo "==> Building for multiple architectures..."
	@GOOS=linux GOARCH=amd64 go build -o dist/$(PLUGIN_NAME)-linux-amd64 .
	@GOOS=linux GOARCH=arm64 go build -o dist/$(PLUGIN_NAME)-linux-arm64 .
	@GOOS=darwin GOARCH=amd64 go build -o dist/$(PLUGIN_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build -o dist/$(PLUGIN_NAME)-darwin-arm64 .
	@GOOS=windows GOARCH=amd64 go build -o dist/$(PLUGIN_NAME)-windows-amd64.exe .
	@echo "==> Multi-architecture build complete!"

clean:
	@echo "==> Cleaning build artifacts..."
	@rm -rf dist/
	@rm -f $(PLUGIN_NAME)
	@rm -f $(PLUGIN_NAME).exe
	@echo "==> Clean complete!"

# Install golangci-lint for local development
install-lint:
	@echo "==> Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2
	@echo "==> golangci-lint installed!"

# Run all quality checks
quality: fmt fmtcheck lint vet
	@echo "==> All quality checks passed!"

# Prepare for release
pre-release: clean quality test ci-build
	@echo "==> Pre-release checks complete!"

.PHONY: build build-local test testacc fmt fmtcheck lint test-compile vendor local-install local-install-modern local-install-dev install-local vet ci-test ci-build clean install-lint quality pre-release
