BASE_PATH=github.com/streamnative/streamnative-mcp-server
VERSION_PATH=main
GIT_VERSION=$(shell git describe --tags --abbrev=0)-SNAPSHOT-$(shell git rev-parse --short HEAD)
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
MKDIR_P = mkdir -p

export GOPRIVATE=github.com/streamnative

.PHONY: all
all: build ;

.PHONY: build
build:
	${MKDIR_P} bin/
	CGO_ENABLED=0 go build -ldflags "\
		-X ${VERSION_PATH}.version=${GIT_VERSION} \
		-X ${VERSION_PATH}.commit=${GIT_COMMIT} \
		-X ${VERSION_PATH}.date=${BUILD_DATE}" \
		-o bin/snmcp cmd/streamnative-mcp-server/main.go
