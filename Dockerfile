FROM alpine:3.21
RUN apk --no-cache add ca-certificates
WORKDIR /server
COPY snmcp /server/snmcp

ENTRYPOINT ["/server/snmcp"]
