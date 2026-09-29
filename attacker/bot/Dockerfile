
# Multi-stage Docker build file for bots

# 1. Build stage
FROM golang:1.27-alpine AS builder
WORKDIR /app

# TODO: Can include tcpdump to show network activity

# Downloading dependencies
COPY go.mod go.sum ./
RUN go mod download

# Compiling bot scrips
COPY bot/ ./bot/
COPY cmd/bot/ ./cmd/bot/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bot ./cmd/bot

# 2. Runtime stage
FROM alpine:latest

# Runs the binary which starts the relevant scripts
RUN adduser -D -H bot

COPY --from=builder /bot /bot
RUN chmod 755 /bot && chown bot:bot /bot
USER bot

ENTRYPOINT ["/bot"]