FROM golang:1.21.3-alpine AS builder

COPY .. /app
WORKDIR /app

RUN go mod download
RUN go mod tidy
RUN go build -o ./bin/auth_server cmd/auth/main.go

FROM alpine:3.18

WORKDIR /root/
COPY --from=builder /app/bin/auth_server .
ADD config/.env .
ENV APP_ENV=prod

CMD ["./auth_server"]
