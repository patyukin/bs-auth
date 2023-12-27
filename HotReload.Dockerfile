FROM golang:1.21.5-alpine3.19

ENV config=docker

WORKDIR /app

COPY . .
COPY .env ./bin/.env
ENV ENV_FILE_PATH=.env

RUN go mod download
RUN go mod tidy
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

EXPOSE 11000
ENTRYPOINT CompileDaemon --build="go build -o bin/auth_server cmd/auth/main.go" --command="./bin/auth_server"
