# Build stage
FROM golang:1.23-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Run stage
FROM ubuntu:latest
ENV DEBIAN_FRONTEND=noninteractive
COPY --from=builder /app/main .
COPY --from=builder /app/config.env .
RUN apt update && apt install -y \
    mysql-client-core-8.0

CMD ["./main"]