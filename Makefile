GIT_VERSION = $(shell git describe --always --dirty --tags 2> /dev/null || echo 'unversioned')

VERSION_FILE = ./version.go
PKG_VERSION ?= $(GIT_VERSION)

setver:
	@sed -i '' 's/\ver = ".*"/ver = "$(PKG_VERSION)"/' $(VERSION_FILE)

test:
	@go test ./...
	@go clean -testcache

version:
	@echo $(GIT_VERSION)

.PHONY: config test version
