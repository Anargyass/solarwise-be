# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /app/solarwise-api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates \
    && adduser -D -g "" app
WORKDIR /app
COPY --from=builder /app/solarwise-api /app/solarwise-api

USER app
EXPOSE 8080
ENTRYPOINT ["/app/solarwise-api"]
