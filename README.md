#  ToDoCLI

  

Приложение для управления списком задач (TODO) через:

  

- Консольный интерфейс (CLI)

- REST API с авторизацией (JWT)

- gRPC-сервер

- Покрытие тестами (unit, integration)

- Развёртывание через Docker Compose

- Swagger-документация

  

---

  

##  Запуск

  

###  Через Docker Compose

  

```bash

docker compose up --build

```

  

API будет доступен по: http://localhost:8080

  

Swagger-документация: http://localhost:8080/docs/index.html

  

gRPC сервер работает на localhost:50051

  

### Через Go

REST-сервер:

```
```bash

go run main.go
```

### CLI:

```
```bash
go run cmd/cli/main.go
```

  

### Тестирование

Запуск всех тестов:

  


```bash
go test ./...
```

#### Возможности  CLI

- Добавление задачи

- Просмотр всех задач

- Поиск по ID

- Обновление

- Удаление

  

### REST API

```

POST /login Получить JWT-токен

GET /api/items Получить все задачи

GET /api/item/{id} Получить задачу по ID

POST /api/item Создать задачу

PUT /api/item/{id} Обновить задачу

DELETE /api/item/{id} Удалить задачу

  

Все API (кроме /login) требуют JWT
```

  

### Swagger

Документация доступна по адресу:

```bash
http://localhost:8080/docs/index.html
```

### Архитектура
```
cmd/ — точки входа CLI и gRPC

internal/

handlers/ — HTTP контроллеры

service/ — бизнес-логика

repository/ — хранилища (InMemory, Postgres)

model/ — структура задачи

proto/ — описание gRPC API

pb/ — сгенерированные protobuf файлы

docs/ — Swagger-документация
```
