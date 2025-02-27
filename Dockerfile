
FROM golang:latest

WORKDIR /app

COPY . .
RUN go mod tidy
RUN go build -o agent cmd/main.go

CMD ["./agent"]
