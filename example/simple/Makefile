APP_NAME=greeting
APP_VERSION=v0.1.0
APP_BUILD=$(shell date -u +%Y-%m-%d.%H:%M)

GO_PKG=github.com/akaahmedkamal/go-cli

GO_LDFLAGS= -X $(GO_PKG)/v1.AppName=$(APP_NAME) \
			-X $(GO_PKG)/v1.AppVersion=$(APP_VERSION) \
			-X $(GO_PKG)/v1.AppBuild=$(APP_BUILD)

build:
	go build -ldflags "$(GO_LDFLAGS)" -o bin/$(APP_NAME) .
.PHONY: build
