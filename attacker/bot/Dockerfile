
# Multi-stage Docker build file for bots

# 1. Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Downloading dependencies
COPY go.mod go.sum ./
RUN go mod download

# Compiling bot scrips
COPY bot/ ./bot/
RUN CGO_ENABLED=0 GOOS=linux go build -o /bot ./bot

# 2. Runtime stage
FROM alpine:latest

RUN adduser -D -H bot
USER bot

# Runs the binary which starts the relevant scripts
COPY --from=builder /bot /bot
ENTRYPOINT ["/bot"]