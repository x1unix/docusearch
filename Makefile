GO ?= go

CFG 				?= config.dev.yml
E2E_CONFIG_FILE 	?= ../config.dev.yml
PKG 				= ./cmd/docusearchd

.PHONY: all
all: build test

.PHONY: gen
gen:
	@$(GO) generate ./internal/services/store/store_synced_test.go

.PHONY: build
build:
	@$(GO) build $(PKG)

.PHONY: run
run:
	@$(GO) run $(PKG) -config $(CFG)

.PHONY: test
test:
	@$(GO) test -v $(shell $(GO) list ./... | grep -v e2e)

.PHONY: e2e
e2e:
	@E2E_CONFIG_FILE=$(E2E_CONFIG_FILE) $(GO) test -v ./e2e/...