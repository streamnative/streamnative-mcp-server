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

# Change ownership
RUN chown -R snmcp:snmcp /server

# Switch to non-root user
USER snmcp

# Expose port if needed (adjust based on your application)
# EXPOSE 8080

ENTRYPOINT ["/server/snmcp"]
