FROM golang:1.24.2-alpine

WORKDIR /app

# Установка зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Генерация Swagger
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4 \
    && swag init -g cmd/server/main.go

# Сборка CLI и сервера
RUN go build -o server ./cmd/server

EXPOSE 8080

CMD ["./server"]
