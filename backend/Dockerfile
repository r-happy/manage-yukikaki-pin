# syntax=docker/dockerfile:1

FROM golang:1.24.3-bookworm AS builder

ENV GOTOOLCHAIN=go1.24.3+auto

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src ./src

RUN CGO_ENABLED=1 GOOS=linux go build -o bin/backend ./src

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        libsqlite3-0 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/bin/backend ./backend
COPY --from=builder /app/src/db ./db

EXPOSE 1323

CMD ["./backend"]
