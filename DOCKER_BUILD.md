# Building Multi-Platform Docker Images

This document explains how to build Docker images that work on both Linux and Mac (including M1/M2/M3 Macs).

## Prerequisites

- Docker 19.03 or newer with buildx support
- Docker buildx plugin enabled

## Quick Start

### Build for Local Platform

To build a Docker image for your current platform:

```bash
make docker-build
```

This creates images with both names for backward compatibility:
- `docker.io/streamnative/mcp-server:latest`
- `docker.io/streamnative/mcp-server:<git-version>`
- `docker.io/streamnative/snmcp:latest` (legacy name)
- `docker.io/streamnative/snmcp:<git-version>` (legacy name)

### Build Multi-Platform Images

To build images for both Linux (amd64) and Mac (arm64):

```bash
make docker-build-multiplatform
```

### Build and Push to Registry

To build multi-platform images and push them to a registry:

```bash
make docker-build-push
```

## Customization

You can customize the build using environment variables:

```bash
# Use custom registry
DOCKER_REGISTRY=myregistry.io make docker-build-push

# Use custom image names
DOCKER_IMAGE=my-mcp-server DOCKER_IMAGE_LEGACY=my-snmcp make docker-build

# Build for specific platforms
DOCKER_PLATFORMS="linux/amd64,linux/arm64,linux/arm/v7" make docker-build-multiplatform
```

## Platform Support

By default, images are built for:
- `linux/amd64` - Intel/AMD 64-bit processors (most Linux servers and older Macs)
- `linux/arm64` - ARM 64-bit processors (Apple Silicon Macs, ARM-based servers)

## Running the Image

You can use either `mcp-server` or `snmcp` image names - they are identical:

### On Linux (amd64/x86_64)

```bash
# Using the new name
docker run --rm -it docker.io/streamnative/mcp-server:latest

# Using the legacy name (backward compatibility)
docker run --rm -it docker.io/streamnative/snmcp:latest
```

### On Mac (Intel)

```bash
# Using the new name
docker run --rm -it docker.io/streamnative/mcp-server:latest

# Using the legacy name (backward compatibility)
docker run --rm -it docker.io/streamnative/snmcp:latest
```

### On Mac (Apple Silicon M1/M2/M3)

```bash
# Using the new name
docker run --rm -it --platform linux/arm64 docker.io/streamnative/mcp-server:latest

# Using the legacy name (backward compatibility)
docker run --rm -it --platform linux/arm64 docker.io/streamnative/snmcp:latest
```

## Image Names

For backward compatibility, all builds create images with two names:
- `streamnative/mcp-server` - The new, preferred name
- `streamnative/snmcp` - The legacy name for backward compatibility

Both names point to the exact same image and can be used interchangeably.

## Troubleshooting

### Enable Docker Buildx

If buildx is not available, enable it:

```bash
# Check if buildx is available
docker buildx version

# Create and use a new builder
docker buildx create --use
```

### Clean Build Environment

If you encounter issues, clean the buildx environment:

```bash
make docker-buildx-clean
make docker-buildx-setup
```

## Build Arguments

The Dockerfile accepts the following build arguments:
- `VERSION` - Version string for the binary
- `COMMIT` - Git commit hash
- `BUILD_DATE` - Build timestamp

These are automatically set by the Makefile based on git information.

## Security

The Docker image:
- Runs as a non-root user (uid: 1000, gid: 1000)
- Uses minimal Alpine Linux base image
- Includes only necessary CA certificates for HTTPS connections
- Has a read-only root filesystem (can be enforced with `--read-only` flag)

## Example: Complete Build and Run

```bash
# Setup buildx (first time only)
make docker-buildx-setup

# Build multi-platform image (creates both mcp-server and snmcp tags)
make docker-build-multiplatform

# Run on current platform (using new name)
docker run --rm -it docker.io/streamnative/mcp-server:latest

# Run on current platform (using legacy name)
docker run --rm -it docker.io/streamnative/snmcp:latest

# Run with custom configuration
docker run --rm -it \
  -v $(pwd)/config.yaml:/server/config.yaml:ro \
  docker.io/streamnative/mcp-server:latest \
  --config /server/config.yaml
``` 