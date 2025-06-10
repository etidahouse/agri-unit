
FROM golang:1.24-alpine AS builder

ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -ldflags="-s -w" -o /app/ingestor .

FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /root/

COPY --from=builder /app/ingestor .

RUN chmod +x /root/ingestor

EXPOSE 8080

ENTRYPOINT ["/root/ingestor"]
