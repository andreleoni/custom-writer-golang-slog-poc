FROM golang:1.22.1

WORKDIR /go/src/app

COPY . .

RUN go mod download
