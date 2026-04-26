APP_NAME=flowgo
MAIN=main.go
BUILD_DIR=build

GO=go
GOFLAGS=

VERSION?=dev

LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Default
all: build

## Build binary
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)

## Run
run:
	$(GO) run $(MAIN)

## Install to system (GOBIN or GOPATH/bin)
install:
	@echo "Installing $(APP_NAME)..."
	$(GO) install $(LDFLAGS)

## Install to custom directory (example: /usr/local/bin)
install-global:
	@echo "Installing to /usr/local/bin (may require sudo)..."
	@mkdir -p /usr/local/bin
	$(GO) build $(LDFLAGS) -o /usr/local/bin/$(APP_NAME) $(MAIN)

## Uninstall
uninstall:
	@echo "Removing $(APP_NAME)..."
	@rm -f $(GOBIN)/$(APP_NAME)
	@rm -f $(GOPATH)/bin/$(APP_NAME)
	@rm -f /usr/local/bin/$(APP_NAME)

## Clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

## Format
fmt:
	$(GO) fmt ./...

## Lint
lint:
	golangci-lint run

## Cross build
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)

	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN)
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN)
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-mac-amd64 $(MAIN)

## Help
help:
	@echo "Available commands:"
	@echo "  make build            - Build binary"
	@echo "  make run              - Run app"
	@echo "  make install          - Install to Go bin"
	@echo "  make install-global   - Install to /usr/local/bin"
	@echo "  make uninstall        - Remove installed binary"
	@echo "  make clean            - Remove build files"
	@echo "  make fmt              - Format code"
	@echo "  make lint             - Run linter"
	@echo "  make build-all        - Cross compile"
