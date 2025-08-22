FROM golang:1.24.3-alpine AS builder
LABEL authors="sakuheinonen"

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -tags reader -a -o main main.go

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/main .

ENTRYPOINT ["./main"]