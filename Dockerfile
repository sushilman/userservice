FROM golang:1.17-alpine3.13 as builder

WORKDIR /go/src/userservice

COPY go.mod .
COPY go.sum .
RUN GO111MODULE=on go mod download
COPY . .
RUN GO111MODULE=on go build -o app main.go

CMD ["./app"]
