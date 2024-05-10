FROM golang:1.21-buster as build

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY main.go /app/

RUN go build -o app main.go

RUN ./app
