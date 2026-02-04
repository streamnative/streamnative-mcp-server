# Copyright 2025 StreamNative
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Multi-stage build for multi-platform support
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files and SDK directories first
COPY go.mod go.sum go.work go.work.sum ./
COPY sdk/ ./sdk/

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build arguments for cross-compilation
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG COMMIT
ARG BUILD_DATE

# Build the binary for the target platform
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags "\
        -X main.version=${VERSION} \
        -X main.commit=${COMMIT} \
        -X main.date=${BUILD_DATE}" \
    -o snmcp cmd/streamnative-mcp-server/main.go

# Final stage - minimal Alpine image
FROM alpine:3.21

# Install CA certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 snmcp && \
    adduser -D -u 1000 -G snmcp snmcp

# Set working directory
WORKDIR /server

# Copy binary from builder
COPY --from=builder /build/snmcp /server/snmcp
COPY --from=builder /build/cmd/snmcp-e2e/testdata/functions/echo.py /server/e2e/functions/echo.py

# Change ownership
RUN chown -R snmcp:snmcp /server

# Switch to non-root user
USER snmcp

# Expose port if needed (adjust based on your application)
# EXPOSE 8080

ENTRYPOINT ["/server/snmcp"]
