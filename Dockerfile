FROM alpine:3.21
RUN apk --no-cache add ca-certificates
WORKDIR /server
COPY streamnative-mcp-server /server/snmcp

ENTRYPOINT ["/server/snmcp"]
