FROM golang:1.12 AS builder

WORKDIR /app
COPY . .

ENV GO111MODULE=on
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o expected-registryhook cmd/expected-registryhook/main.go

ENTRYPOINT ["/app/expected-registryhook"]
