GIT_VERSION = $(shell git describe --always --dirty --tags 2> /dev/null || echo 'unversioned')

VERSION_FILE = ./pkg/version/version.go
PKG_VERSION ?= $(GIT_VERSION)

config:
	@sed -i '' 's/\Version = ".*"/Version = "$(PKG_VERSION)"/' $(VERSION_FILE)

test:	## Tests packages
	@go test ./...

version:	## Shows the current git version
	@echo $(GIT_VERSION)

help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: config test version help
