# Spotify Shuffle - Go Edition
# Makefile for building and packaging

.PHONY: help build build-all clean test deps dev install package-all release

# Variables
APP_NAME := spotify-shuffle
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR := build
DIST_DIR := dist
LDFLAGS := -s -w -X main.version=$(VERSION)

# Default target
help: ## Show this help message
	@echo "Spotify Shuffle - Go Edition"
	@echo "============================="
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Install dependencies
	go mod download
	go mod tidy

test: ## Run tests
	go test -v ./...

dev: ## Build with race detection for development
	go build -race -o $(BUILD_DIR)/$(APP_NAME) .

build: ## Build for current platform
	mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) .

build-all: ## Build for all platforms
	mkdir -p $(BUILD_DIR)
	
	# Linux
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 .
	
	# macOS
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-macos-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-macos-arm64 .
	
	# Windows
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe .

install: build ## Install to system (requires sudo on Linux/macOS)
	@if [ "$(shell uname)" = "Darwin" ] || [ "$(shell uname)" = "Linux" ]; then \
		sudo cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/; \
		echo "Installed to /usr/local/bin/$(APP_NAME)"; \
	elif [ "$(OS)" = "Windows_NT" ]; then \
		echo "Please manually copy $(BUILD_DIR)/$(APP_NAME).exe to a directory in your PATH"; \
	fi

clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR) $(DIST_DIR)

# Packaging targets (requires platform-specific tools)
package-deb: ## Create DEB package (Linux only)
	@if [ "$(shell uname)" != "Linux" ]; then \
		echo "DEB packaging only supported on Linux"; \
		exit 1; \
	fi
	mkdir -p $(DIST_DIR)/deb/DEBIAN
	mkdir -p $(DIST_DIR)/deb/usr/local/bin
	cp $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(DIST_DIR)/deb/usr/local/bin/$(APP_NAME)
	chmod +x $(DIST_DIR)/deb/usr/local/bin/$(APP_NAME)
	echo "Package: $(APP_NAME)" > $(DIST_DIR)/deb/DEBIAN/control
	echo "Version: $(VERSION)" >> $(DIST_DIR)/deb/DEBIAN/control
	echo "Section: sound" >> $(DIST_DIR)/deb/DEBIAN/control
	echo "Priority: optional" >> $(DIST_DIR)/deb/DEBIAN/control
	echo "Architecture: amd64" >> $(DIST_DIR)/deb/DEBIAN/control
	echo "Maintainer: Spotify Shuffle <noreply@example.com>" >> $(DIST_DIR)/deb/DEBIAN/control
	echo "Description: CLI tool for managing Spotify playlists" >> $(DIST_DIR)/deb/DEBIAN/control
	echo " A fast, cross-platform CLI tool for managing your Spotify playlists." >> $(DIST_DIR)/deb/DEBIAN/control
	dpkg-deb --build $(DIST_DIR)/deb $(DIST_DIR)/$(APP_NAME)-$(VERSION)-amd64.deb

package-dmg: ## Create DMG package (macOS only)
	@if [ "$(shell uname)" != "Darwin" ]; then \
		echo "DMG packaging only supported on macOS"; \
		exit 1; \
	fi
	mkdir -p $(DIST_DIR)/dmg
	cp $(BUILD_DIR)/$(APP_NAME)-macos-amd64 $(DIST_DIR)/dmg/$(APP_NAME)
	chmod +x $(DIST_DIR)/dmg/$(APP_NAME)
	hdiutil create -volname "Spotify Shuffle" -srcfolder $(DIST_DIR)/dmg -ov -format UDZO $(DIST_DIR)/$(APP_NAME)-$(VERSION)-macos.dmg

package-all: build-all ## Build and package for all platforms (requires platform-specific tools)
	@echo "Building packages for all platforms..."
	@echo "Note: Some packages may fail if platform-specific tools are not available"
	-$(MAKE) package-deb
	-$(MAKE) package-dmg

release: ## Tag and push a new release
	@echo "Current version: $(VERSION)"
	@read -p "Enter new version (e.g., v1.0.0): " NEW_VERSION; \
	git tag -a $$NEW_VERSION -m "Release $$NEW_VERSION"; \
	git push origin $$NEW_VERSION; \
	echo "Tagged and pushed $$NEW_VERSION"

# Development helpers
fmt: ## Format code
	go fmt ./...

lint: ## Run linter (requires golangci-lint)
	golangci-lint run

vet: ## Run go vet
	go vet ./...

check: fmt vet lint test ## Run all checks