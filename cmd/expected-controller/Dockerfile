FROM golang:latest AS builder
ENV GO111MODULE=on

WORKDIR /app

COPY cmd ./cmd/
COPY pkg ./pkg/

COPY vendor ./vendor/
COPY go.mod .
COPY go.sum .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -o /go/bin/expected-controller cmd/expected-controller/main.go

FROM alpine

COPY --from=builder /go/bin/expected-controller /go/bin/

ENTRYPOINT ["/go/bin/expected-controller"]
