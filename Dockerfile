FROM golang:1.20

WORKDIR /app

COPY . .

RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init
RUN go build -o todocli main.go

EXPOSE 8080

CMD ["./todocli"]
