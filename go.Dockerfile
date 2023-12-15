FROM golang:1.21.3-alpine AS builder

COPY . /app
WORKDIR /app

RUN go mod download
RUN go mod tidy
RUN go build -o bin/auth_server cmd/auth/main.go

FROM alpine:3.18

WORKDIR /root/
COPY --from=builder /app/bin/auth_server .
COPY .env .
ENV ENV_FILE_PATH=.env

CMD ["./auth_server"]

