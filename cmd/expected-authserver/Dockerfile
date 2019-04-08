FROM golang:latest AS builder
ENV GO111MODULE=on

WORKDIR /app

COPY cmd ./cmd/
COPY pkg ./pkg/

COPY vendor ./vendor/
COPY go.mod .
COPY go.sum .

RUN GOOS=linux GOARCH=amd64 go build -mod vendor -o /go/bin/expected-authserver cmd/expected-authserver/main.go

FROM alpine

COPY --from=builder /go/bin/expected-authserver /go/bin/

ENTRYPOINT ["/go/bin/expected-authserver"]
