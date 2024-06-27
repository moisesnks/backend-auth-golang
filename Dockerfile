# Build stage
FROM golang:1.21.0-alpine3.17 AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod download

COPY *.go .
RUN go build -o app

# Exponer el puerto 8081
EXPOSE 8081

ENTRYPOINT [ "./app" ]