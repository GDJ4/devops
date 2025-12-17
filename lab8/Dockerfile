# syntax=docker/dockerfile:1

FROM golang:1.21-bookworm AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o todo ./cmd/server

FROM alpine:3.19
RUN adduser -D -H -s /sbin/nologin app
WORKDIR /srv

COPY --from=builder /app/todo ./todo

ENV PORT=8080
EXPOSE 8080
USER app

ENTRYPOINT ["./todo"]
