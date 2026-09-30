BINARY_NAME := recall
BUILD_DIR := bin
MODULE := github.com/managedkaos/recall
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64
PREFIX ?= $(HOME)/.local
BINDIR ?= $(PREFIX)/bin
MANDIR ?= $(PREFIX)/share/man

# Local version from the nearest matching Git tag; environment/CLI overrides win.
VERSION ?= $(shell git describe --tags --dirty=-local --match 'v[0-9]*' --match 'V[0-9]*' --match '[0-9]*' 2>/dev/null || echo unknown)

# Build metadata (overridable via environment, e.g. in CI)
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo unknown)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
BUILD_ENV  ?= local ($(shell uname -s))

LDFLAGS := -ldflags "-X '$(MODULE)/cmd.Version=$(VERSION)' -X '$(MODULE)/cmd.GitBranch=$(GIT_BRANCH)' -X '$(MODULE)/cmd.BuildDate=$(BUILD_DATE)' -X '$(MODULE)/cmd.BuildEnvironment=$(BUILD_ENV)'"

.PHONY: help
help: ## Display available targets
	@awk 'BEGIN {FS = ":.*## "}; /^[a-zA-Z0-9_-]+:.*## / {printf "\033[36m%-28s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: all
all: clean test build ## Clean artifacts, run tests, and build the binary

.PHONY: build
build: ## Build the binary for the current platform
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

.PHONY: snapshot
snapshot: ## Build a local release snapshot with GoReleaser (no publish)
	@command -v goreleaser >/dev/null 2>&1 || { echo "goreleaser not found; install from https://goreleaser.com/install/"; exit 1; }
	goreleaser release --snapshot --clean

.PHONY: man
man: ## Generate man pages in ./man using Cobra
	go run ./tools/manpages

.PHONY: clean
clean: ## Remove build artifacts
	rm -rvf $(BUILD_DIR) dist man

.PHONY: install
install: man ## Build a snapshot and install the native binary and man pages (Linux/macOS)
	@case "$$(uname -s)" in Linux|Darwin) ;; *) echo "install supports Linux and macOS only" >&2; exit 1 ;; esac
	$(MAKE) snapshot
	@os_name=$$(uname -s | tr '[:upper:]' '[:lower:]'); \
	case "$$(uname -m)" in x86_64|amd64) arch_name=amd64 ;; arm64|aarch64) arch_name=arm64 ;; *) echo "unsupported architecture" >&2; exit 1 ;; esac; \
	mkdir -p "$(DESTDIR)$(BINDIR)" "$(DESTDIR)$(MANDIR)/man1" && \
	install -m 0755 dist/$(BINARY_NAME)_$${os_name}_$${arch_name}*/$(BINARY_NAME) "$(DESTDIR)$(BINDIR)/$(BINARY_NAME)" && \
	install -m 0644 man/*.1 "$(DESTDIR)$(MANDIR)/man1/"
