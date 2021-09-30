FROM golang:1.17.1

ENV GO111MODULE=on

WORKDIR /alfred

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build

ENTRYPOINT ./alfred