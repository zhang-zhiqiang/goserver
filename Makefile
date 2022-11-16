# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# Build all by default, even if it's not first
.DEFAULT_GOAL := all

.PHONY: all
#all: tidy format lint cover build
all: tidy format lint cover build

# ==============================================================================
# Build options

ROOT_PACKAGE=github.com/marmotedu/goserver
VERSION_PACKAGE=github.com/marmotedu/goserver/pkg/version

# ==============================================================================
# Includes

include scripts/make-rules/common.mk # make sure include common.mk at the first include line
include scripts/make-rules/golang.mk
include scripts/make-rules/tools.mk

# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  DEBUG            Whether to generate debug symbols. Default is 0.
  BINS             The binaries to build. Default is all of cmd.
                   This option is available when using: make build/build.multiarch
                   Example: make build BINS="goserver test"
  VERSION          The version information compiled into binaries.
                   The default is obtained from gsemver or git.
  V                Set to 1 enable verbose build. Default is 0.
endef
export USAGE_OPTIONS

## --------------------------------------
## Generate / Manifests
## --------------------------------------

##@ generate:

.PHONY: ca
ca: ## Generate CA files for goserver.
	@mkdir -p $(OUTPUT_DIR)/cert
	@openssl req -new -nodes -x509 -out $(OUTPUT_DIR)/cert/goserver.pem -keyout $(OUTPUT_DIR)/cert/goserver-key.pem \
		-days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=127.0.0.1/emailAddress=xxxxx@qq.com"


## --------------------------------------
## Binaries
## --------------------------------------

##@ build:

.PHONY: build
build: ## Build source code for host platform.
	@$(MAKE) go.build

## --------------------------------------
## Cleanup
## --------------------------------------

##@ clean:

.PHONY: clean
clean: ## Remove all files that are created by building and generaters.
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)


## --------------------------------------
## Lint / Verification
## --------------------------------------

##@ lint and verify:

.PHONY: lint
lint: ## Check syntax and styling of go sources.
	@$(MAKE) go.lint


## --------------------------------------
## Testing
## --------------------------------------

##@ test:

.PHONY: test 
test: ## Run unit test.
	@$(MAKE) go.test

.PHONY: cover 
cover: ## Run unit test and get test coverage.
	@$(MAKE) go.test.cover


## --------------------------------------
## Hack / Tools
## --------------------------------------

##@ hack/tools:

.PHONY: format
format: tools.verify.goimports ## Gofmt (reformat) package sources (exclude vendor dir if existed). 
	@echo "===========> Formating codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(GO) mod edit -fmt

.PHONY: deps
deps: ## Install required dependencies, such as tools.
	@$(MAKE) tools.install

.PHONY: tidy
tidy:
	@$(GO) mod tidy

.PHONY: help
help: Makefile ## Show this help info.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<TARGETS> <OPTIONS>\033[0m\n\n\033[35mTargets:\033[0m\n"} /^[0-9A-Za-z._-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^\$$\([0-9A-Za-z_-]+\):.*?##/ { gsub("_","-", $$1); printf "  \033[36m%-45s\033[0m %s\n", tolower(substr($$1, 3, length($$1)-7)), $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' Makefile #$(MAKEFILE_LIST)
	@echo -e "$$USAGE_OPTIONS"
