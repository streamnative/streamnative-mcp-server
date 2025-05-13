FROM golang:1.24.3-alpine3.21 AS builder
ARG VERSION="snapshot"

WORKDIR /build

RUN --mount=type=cache,target=/var/cache/apk \
  apk add git

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=bind,target=. \
  CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=${VERSION} -X main.commit=$(git rev-parse HEAD) -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o /bin/snmcp cmd/streamnative-mcp-server/main.go

FROM alpine:3.21
WORKDIR /server
COPY --from=builder /bin/snmcp .

ENTRYPOINT ["/server/snmcp"]
