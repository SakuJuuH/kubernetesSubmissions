FROM golang:1.24.3-alpine AS builder
LABEL authors="sakuheinonen"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/main .

ENTRYPOINT ["./main"]
