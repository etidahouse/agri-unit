FROM golang:1.23.2-alpine3.19

RUN apk add --no-cache bash git gcc musl-dev make curl wget ca-certificates && \
    update-ca-certificates
