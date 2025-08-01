#  ToDoCLI

Приложение для управления задачами (TODO-list), поддерживает:

- ✅ REST API с JWT авторизацией  
- ✅ CLI-клиент, который взаимодействует с сервером  
- ✅ gRPC-сервер  
- ✅ Swagger-документация  
- ✅ Docker Compose для развёртывания  
- ✅ Unit и интеграционные тесты  

---

##  Быстрый старт

### Через Docker Compose

```bash
docker compose up --build
```

 После запуска:

- **REST API**: [http://localhost:8080](http://localhost:8080)  
- **Swagger**: [http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)  
- **gRPC-сервер**: `localhost:50051`

> Убедитесь, что у вас настроен `.env` файл с переменными (`API_USER`, `API_PASS`, `JWT_SECRET`, и т.д.)

---

### Через Go

#### Запуск REST сервера:

```bash
go run cmd/server/main.go
```

#### Запуск CLI:

```bash
go run cmd/cli/main.go
```

> CLI сначала выполняет авторизацию (`/login`), затем позволяет управлять задачами через REST API.

---

##  Пример .env файла

```env
API_USER=test
API_PASS=test
JWT_SECRET=supersecretjwtkey
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/tododb?sslmode=disable
REDIS_ADDR=localhost:6379
```

---

##  Возможности CLI

-  Авторизация по логину и паролю  
-  Добавление задачи  
-  Просмотр всех задач  
-  Поиск по ID  
-  Обновление задачи  
-  Удаление задачи

---

##  REST API

```
POST   /login           — Получить JWT-токен
GET    /api/items       — Получить все задачи
GET    /api/item/:id    — Получить задачу по ID
POST   /api/item        — Создать задачу
PUT    /api/item/:id    — Обновить задачу
DELETE /api/item/:id    — Удалить задачу
```

⚠️ Все запросы к `/api/*` требуют JWT токен в заголовке:  
```
Authorization: Bearer <token>
```

---

##  Swagger

Документация доступна по адресу:  
 [http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)

---

##  Тестирование

Запуск всех тестов:

```bash
go test ./...
```

---

## Архитектура проекта

```text
cmd/
  cli/         — CLI клиент
  server/      — Точка входа для REST сервера

internal/
  handlers/    — HTTP-обработчики
  service/     — Бизнес-логика
  repository/  — Хранилища: InMemory, PostgreSQL
  model/       — Определение сущностей (Task)
  handlers/middleware/ — JWT middleware

proto/         — gRPC описание сервиса
pb/            — Сгенерированные файлы Protobuf
docs/          — Swagger-документация
```

---

## Примечания

- Проект покрыт юнит-тестами (`service`, `repository`) и интеграционными тестами (`handlers`).  
- CLI и сервер используют общий код и взаимодействуют по REST.
