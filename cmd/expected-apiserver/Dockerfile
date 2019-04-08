FROM golang:1.12

WORKDIR /app
COPY . .

ENV GO111MODULE=on
RUN go mod download
RUN go build -o expected-apiserver cmd/expected-apiserver/main.go

CMD ["/app/expected-apiserver"]
