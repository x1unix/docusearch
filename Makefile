GO ?= go

CFG = config.dev.yml
PKG = ./cmd/docusearchd

.PHONY: all
all: build

.PHONY: build
build:
	@go build $(PKG)

.PHONY: run
run:
	@go run $(PKG) -config $(CFG)
