FROM golang:latest

WORKDIR /bot

COPY . .

RUN go mod tidy

CMD go run main.go