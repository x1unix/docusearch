# DocuSearch

Simple document upload and search service.

## Goal

Create a simple service that allows to upload text files and search which text files contain a specified word.

## Usage

See [Quick start quite](#quick-start-guide) below to find out how to start service.

API documentation is available as [Swagger spec file](swagger.yml). Test files for upload available in [e2e/testdata](e2e/testdata) directory.

**Important note:** by default, [file indexing engine](internal/services/search/indexer.go) omits common verbs and articles in English language to reduce index size.

You can control this behavior by changing `ignore_common_words` parameter in config file.

## How To Run

### Prerequisites

* Go 1.17+
* GNU Make (BSD also might work)
* Redis (or docker-compose)

### Quick start guide
 
Ensure that you have [docker-compose](https://docs.docker.com/compose/) and latest [Go](golang.org/dl/) versions installed.

1. Clone this repo
2. Create containers with dependencies using `docker-compose up -d` (only once)
3. Start containers with `docker-compose start`
4. Run project using `make run`

Service HTTP listen address is specified in default development config (see [config.dev.yml](config.dev.yml))

## Testing

### Unit tests

Use `make test` to run unit tests.

Mocks can be updated using `make gen` command.

### End-to-end

End-to-end tests are located at [e2e](e2e) directory.

Use `make e2e` to run end-to-end tests. 

**Warning:** E2E tests require *Redis*. Please start docker-compose containers using `docker-compose start` before running end-to-end tests.

