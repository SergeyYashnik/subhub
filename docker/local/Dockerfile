FROM golang:1.26-alpine

RUN apk add --no-cache git gcc musl-dev

RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD swag init -g cmd/api/main.go && air -c .air.toml