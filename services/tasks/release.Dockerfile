
FROM golang:1.24-alpine AS builder

ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -ldflags="-s -w" -o /app/tasks .

FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /root/

COPY --from=builder /app/tasks .

RUN chmod +x /root/tasks

ENTRYPOINT ["/root/tasks"]
