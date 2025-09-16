
FROM golang:1.24.3 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bot ./cmd

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

COPY --from=builder /app/bot .

CMD ["./bot"]