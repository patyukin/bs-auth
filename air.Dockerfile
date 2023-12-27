FROM golang:1.21.5

ENV ENV_FILE_PATH=.env

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go get github.com/cosmtrek/air
RUN go install github.com/cosmtrek/air@latest


ENTRYPOINT ["air"]
