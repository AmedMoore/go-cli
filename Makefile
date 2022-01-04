APP_NAME=go-cli
APP_VERSION=v0.1.0
APP_BUILD=$(shell date -u +%Y-%m-%d.%H:%M)

GO_LDFLAGS= -X github.com/akaahmedkamal/go-cli/v1.AppName=$(APP_NAME) \
			-X github.com/akaahmedkamal/go-cli/v1.AppVersion=$(APP_VERSION) \
			-X github.com/akaahmedkamal/go-cli/v1.AppBuild=$(APP_BUILD)

build:
	go build -ldflags "$(GO_LDFLAGS)" -o bin/$(APP_NAME) .
.PHONY: build
