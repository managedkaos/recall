BINARY_NAME := recall
BUILD_DIR := bin
MODULE := github.com/mjenkins/recall

PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Default: build for current platform
.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) .

# Build for all target platforms
.PHONY: build-all
build-all:
	@mkdir -p $(BUILD_DIR)
	@$(foreach platform,$(PLATFORMS),\
		$(eval OS := $(word 1,$(subst /, ,$(platform))))\
		$(eval ARCH := $(word 2,$(subst /, ,$(platform))))\
		$(eval EXT := $(if $(filter windows,$(OS)),.exe,))\
		echo "Building $(BINARY_NAME)-$(OS)-$(ARCH)$(EXT)..." && \
		CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME)-$(OS)-$(ARCH)$(EXT) . && \
	) true

# Run all tests
.PHONY: test
test:
	go test ./...

# Remove build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
