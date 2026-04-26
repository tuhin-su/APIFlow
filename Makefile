APP_NAME=flowgo
MAIN=main.go
BUILD_DIR=build

GO=go
VERSION?=dev
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

all: build

build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)

run:
	$(GO) run $(MAIN) -l config.json

install:
	$(GO) install $(LDFLAGS)

install-global:
	$(GO) build $(LDFLAGS) -o /usr/local/bin/$(APP_NAME) $(MAIN)

uninstall:
	rm -f $(GOBIN)/$(APP_NAME)
	rm -f $(GOPATH)/bin/$(APP_NAME)
	rm -f /usr/local/bin/$(APP_NAME)

clean:
	rm -rf $(BUILD_DIR)

fmt:
	$(GO) fmt ./...

build-all:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(APP_NAME)-linux $(MAIN)
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(APP_NAME).exe $(MAIN)
