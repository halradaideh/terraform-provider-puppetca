TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=puppetca
PLUGIN_NAME=terraform-provider-puppetca
PLUGIN_VERSION=1.0.0

default: build

build: fmtcheck
	go install

# Install to the legacy plugin path for local testing/development
local-install: build
	@echo "Installing provider to ~/.terraform.d/plugins/puppetca_v$(PLUGIN_VERSION)/$(PLUGIN_NAME)"
	@mkdir -p ~/.terraform.d/plugins/puppetca_v$(PLUGIN_VERSION)
	@cp $(PLUGIN_NAME) ~/.terraform.d/plugins/puppetca_v$(PLUGIN_VERSION)/$(PLUGIN_NAME)
	@echo "Local install complete! Legacy plugin path: ~/.terraform.d/plugins/puppetca_v$(PLUGIN_VERSION)/$(PLUGIN_NAME)"

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v -parallel 20 $(TESTARGS) -timeout 120m

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)

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

.PHONY: build test testacc fmt fmtcheck lint test-compile vendor
