BINARY_NAME := recall
BUILD_DIR := bin
MODULE := github.com/managedkaos/recall
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Version extraction from version.yml
VERSION_MAJOR := $(shell grep '^major:' version.yml | awk '{print $$2}')
VERSION_MINOR := $(shell grep '^minor:' version.yml | awk '{print $$2}')
VERSION_PATCH := $(shell grep '^patch:' version.yml | awk '{print $$2}')
LDFLAGS := -ldflags "-X $(MODULE)/cmd.Major=$(VERSION_MAJOR) -X $(MODULE)/cmd.Minor=$(VERSION_MINOR) -X $(MODULE)/cmd.Patch=$(VERSION_PATCH)"

.PHONY: help
help:
	@echo coming soon

.PHONY: all
all: clean test build

.PHONY: build
build: ## Default: build for current platform
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

.PHONY: build-all
build-all: ## Build for all target platforms
	@mkdir -p $(BUILD_DIR)
	@$(foreach platform,$(PLATFORMS),\
		$(eval OS := $(word 1,$(subst /, ,$(platform))))\
		$(eval ARCH := $(word 2,$(subst /, ,$(platform))))\
		$(eval EXT := $(if $(filter windows,$(OS)),.exe,))\
		echo "Building $(BINARY_NAME)-$(OS)-$(ARCH)$(EXT)..." && \
		CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(OS)-$(ARCH)$(EXT) . && \
	) true

.PHONY: test
test: ## Run all tests
	go test -v ./...

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)
