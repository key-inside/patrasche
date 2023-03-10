GIT_VERSION = $(shell git describe --always --dirty --tags 2> /dev/null || echo 'unversioned')

VERSION_FILE = ./version.go
PKG_VERSION ?= $(GIT_VERSION)

setver:
	@sed -i '' 's/\ver = ".*"/ver = "$(PKG_VERSION)"/' $(VERSION_FILE)

test:
	@go clean -testcache
	@go test ./...

version:
	@echo $(GIT_VERSION)

.PHONY: config test version
