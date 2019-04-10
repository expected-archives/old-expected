FROM golang:1.12

WORKDIR /app
COPY . .

ENV GO111MODULE=on
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o expected-scheduler cmd/expected-scheduler/main.go

ENTRYPOINT ["/app/expected-scheduler"]
