# Stage 1: Build
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN go mod tidy