# REST API GOLANG

REST API для управления пользователями и заказами на Go + PostgreSQL.

[![Go](https://img.shields.io/badge/Go-1.24-blue)](https://go.dev/) [![Gin](https://img.shields.io/badge/Gin-1.10-green)](https://gin-gonic.com/) [![GORM](https://img.shields.io/badge/GORM-1.26-orange)](https://gorm.io/) [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue)](https://www.postgresql.org/)

---

## 🚀 Быстрый старт

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/Solarmatte/rest-api-users-orders-go.git
   cd rest-api-users-orders-go-main
   ```
2. Скопируйте переменные окружения:
   ```bash
   cp .env.example .env
   ```
3. Запустите приложение:
   ```bash
   docker-compose up --build
   ```
4. API будет доступен на [http://localhost:8080](http://localhost:8080)

---

## 📚 Документация API

Swagger: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## 🏗️ Структура проекта

```
├── cmd/               # main.go — точка входа
├── docs/              # Swagger (авто-сгенерировано)
├── internal/          # бизнес-логика и HTTP-слой
│   ├── config/        # конфиг и .env
│   ├── bootstrap/     # инициализация БД, миграции
│   ├── handlers/      # HTTP-контроллеры (Gin)
│   ├── middleware/    # JWT, логирование, Recovery
│   ├── models/        # GORM-модели (users, orders)
│   ├── repositories/  # CRUD-репозитории
│   ├── router/        # маршрутизация и Swagger
│   ├── services/      # бизнес-логика
│   └── utils/         # утилиты (JWT и др.)
├── migrations/        # SQL-скрипты
├── tests/             # unit & integration тесты
├── .env               # переменные окружения
├── docker-compose.yml # Docker Compose (Postgres + app)
├── Dockerfile         # Docker-образ
├── go.mod / go.sum    # зависимости
└── README.md
```

---

## ⚙️ Переменные окружения

Пример в `.env.example`. Основные:

| Переменная         | Описание                |
|--------------------|------------------------|
| DB_HOST            | Хост PostgreSQL        |
| DB_PORT            | Порт PostgreSQL        |
| DB_USER            | Пользователь БД        |
| DB_PASSWORD        | Пароль БД              |
| DB_NAME            | Имя БД                 |
| JWT_SECRET         | Секрет для JWT         |
| APP_PORT           | Порт приложения        |

---

## 🛠️ Базовые команды

| Операция                | Команда                                  |
|-------------------------|------------------------------------------|
| Сборка                  | `go build -o main ./cmd`                 |
| Локальный запуск        | `go run ./cmd`                           |
| Тесты                   | `go test ./...`                          |
| Docker Compose (run)    | `docker-compose up --build`              |
| Docker Compose (stop)   | `docker-compose down`                    |
| Swagger обновить        | `swag init -g cmd/main.go -o docs ...`   |

---

## 📎 Ссылки
- [Swagger UI](http://localhost:8080/swagger/index.html)
- [Документация Gin](https://gin-gonic.com/docs/)
- [Документация GORM](https://gorm.io/docs/)
- [PostgreSQL](https://www.postgresql.org/)

