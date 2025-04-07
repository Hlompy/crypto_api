# Dockerfile
FROM golang:1.24.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o crypto_api ./main.go

RUN sleep 10

EXPOSE 8080

CMD ["./crypto_api"]