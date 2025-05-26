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
	@GOGC=30 golangci-lint run ./$(PKG_NAME)

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
	go vet $<

.PHONY: build build-local test testacc fmt fmtcheck lint test-compile vendor local-install local-install-modern local-install-dev install-local
