BASE_PATH=github.com/streamnative/streamnative-mcp-server
VERSION_PATH=main
GIT_VERSION=$(shell git describe --tags --abbrev=0)-SNAPSHOT-$(shell git rev-parse --short HEAD)
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
MKDIR_P = mkdir -p

# Docker configuration
DOCKER_REGISTRY ?= docker.io
DOCKER_IMAGE ?= streamnative/mcp-server
DOCKER_IMAGE_LEGACY ?= streamnative/snmcp
DOCKER_TAG ?= $(GIT_VERSION)
DOCKER_PLATFORMS ?= linux/amd64,linux/arm64

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

# Build Docker image for local platform with both names
.PHONY: docker-build
docker-build:
	docker build \
		--build-arg VERSION=$(GIT_VERSION) \
		--build-arg COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_LEGACY):$(DOCKER_TAG) \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_LEGACY):latest \
		.

# Build multi-platform Docker image and push to registry with both names
.PHONY: docker-build-push
docker-build-push: docker-buildx-setup
	docker buildx build \
		--platform $(DOCKER_PLATFORMS) \
		--build-arg VERSION=$(GIT_VERSION) \
		--build-arg COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_LEGACY):$(DOCKER_TAG) \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_LEGACY):latest \
		--push \
		.

# Build multi-platform Docker image without pushing (for testing) with both names
.PHONY: docker-build-multiplatform
docker-build-multiplatform: docker-buildx-setup
	docker buildx build \
		--platform $(DOCKER_PLATFORMS) \
		--build-arg VERSION=$(GIT_VERSION) \
		--build-arg COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_LEGACY):$(DOCKER_TAG) \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_LEGACY):latest \
		--load \
		.

# Setup Docker buildx for multi-platform builds
.PHONY: docker-buildx-setup
docker-buildx-setup:
	@if ! docker buildx ls | grep -q mcp-builder; then \
		docker buildx create --name mcp-builder --use; \
		docker buildx inspect --bootstrap; \
	else \
		docker buildx use mcp-builder; \
	fi

# Clean Docker buildx builder
.PHONY: docker-buildx-clean
docker-buildx-clean:
	-docker buildx rm mcp-builder

.PHONY: license-check
license-check:
	license-eye header check

# go install github.com/apache/skywalking-eyes/cmd/license-eye@latest
.PHONY: license-fix
license-fix:
	license-eye header fix

# E2E Testing targets
.PHONY: test-e2e-pulsar-start
test-e2e-pulsar-start:
	docker-compose -f test/docker-compose-pulsar.yml up -d
	@echo "Waiting for Pulsar to be ready..."
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		if curl -s http://localhost:8080/admin/v2/clusters > /dev/null 2>&1; then \
			echo "Pulsar is ready!"; \
			exit 0; \
		fi; \
		echo "Waiting for Pulsar... ($$i/10)"; \
		sleep 3; \
	done
	@echo "Error: Pulsar failed to become ready"; \
	exit 1

.PHONY: test-e2e-pulsar-stop
test-e2e-pulsar-stop:
	docker-compose -f test/docker-compose-pulsar.yml down

.PHONY: test-e2e-pulsar-logs
test-e2e-pulsar-logs:
	docker-compose -f test/docker-compose-pulsar.yml logs

.PHONY: test-e2e
test-e2e:
	go test -tags=e2e -v ./pkg/mcp/e2e/...

.PHONY: test-e2e-race
test-e2e-race:
	go test -race -tags=e2e -v ./pkg/mcp/e2e/...

.PHONY: test-e2e-run
test-e2e-run: test-e2e-pulsar-start
	@echo "Running E2E tests..."
	@$(MAKE) test-e2e; \
	result=$$?; \
	$(MAKE) test-e2e-pulsar-stop; \
	exit $$result
