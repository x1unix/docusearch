GO ?= go

CFG = config.dev.yml
PKG = ./cmd/docusearchd

.PHONY: all
all: build

.PHONY: gen
gen:
	@go generate ./internal/services/store/store_synced_test.go

.PHONY: build
build:
	@go build $(PKG)

.PHONY: run
run:
	@go run $(PKG) -config $(CFG)

