GIT_VERSION = $(shell git describe --always --dirty --tags 2> /dev/null || echo 'unversioned')

VERSION_FILE = ./pkg/version/version.go
PKG_VERSION ?= $(GIT_VERSION)

config:
	@sed -i '' 's/\Version = ".*"/Version = "$(PKG_VERSION)"/' $(VERSION_FILE)

test:
	@go test ./...

version:
	@echo $(GIT_VERSION)

.PHONY: config test version
