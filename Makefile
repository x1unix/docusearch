GO ?= go

CFG 				?= config.dev.yml
E2E_CONFIG_FILE 	?= ../$(CFG)
PKG 				= ./cmd/docusearchd

.PHONY: all
all: build test

.PHONY: gen
gen:
	@go generate ./internal/services/store/store_synced_test.go

.PHONY: build
build:
	@go build $(PKG)

.PHONY: run
run:
	@go run $(PKG) -config $(CFG)

.PHONY: test
test:
	@go test -v ./...

.PHONY: e2e
e2e:
	@E2E_CONFIG_FILE=$(E2E_CONFIG_FILE) go test -v ./e2e/...