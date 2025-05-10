
# Kvant Task API

REST API для управления пользователями и их заказами на Go + PostgreSQL.

![Go](https://img.shields.io/badge/Go-1.24-blue) ![Gin](https://img.shields.io/badge/Gin-1.10-green) ![GORM](https://img.shields.io/badge/GORM-1.26-orange) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue)


---

## 🚀 Быстрый старт

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/Solarmatte/rest-api-users-orders-go.git
   cd rest-api-users-orders-go
   ````

2. Отредактируйте под себя образец окружения .env:

   ```bash
   # Сервер
    SERVER_ADDRESS=:8080

    # PostgreSQL
    POSTGRES_HOST=localhost
    POSTGRES_PORT=5432
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=your_password
    POSTGRES_DB=kvant_db
    POSTGRES_SSLMODE=disable

    # JWT
    JWT_SECRET=your_jwt_secret
   ```

3. Запустите PostgreSQL и приложение через Docker Compose:

   ```bash
   docker-compose up -d
   ```

4. Перейдите в Swagger UI:

   ```
   http://localhost:8080/swagger/index.html
   ```

---

## 🏗️ Структура проекта

```
Project/
├── cmd/               # точка входа приложения (main.go)
├── docs/              # Swagger-документация (авто-сгенерированные файлы)  
├── internal/          # основная бизнес-логика и HTTP-слой
│   ├── config/        # загрузка и валидация .env / переменных окружения
│   ├── bootstrap/     # инициализация БД и выполнение миграций
│   ├── handlers/      # HTTP-контроллеры (Gin-эндпоинты)
│   ├── middleware/    # JWT-авторизация, логирование, Recovery
│   ├── models/        # GORM-модели (таблицы users, orders)
│   ├── repositories/  # CRUD-репозитории для работы с БД
│   ├── router/        # настройка маршрутизации Gin и Swagger-UI
│   ├── services/      # бизнес-логика (регистрация, заказы, JWT)
│   └── utils/         # утилиты (JWT-токены и пр.)
├── migrations/        # SQL-скрипты создания таблиц  
├── tests/             # unit & integration тесты  
├── .env               # переменные окружения  
├── .gitignore         # игнорируемые файлы  
├── docker-compose.yml # Docker Compose для Postgres + app  
├── Dockerfile         # сборка Docker-образа  
├── go.mod             # зависимости и модульный путь  
└── go.sum             # контрольные суммы зависимостей  



```

## 🧪 Тестирование

* Unit-tests сервисов и репозиториев лежат в `tests/`.
* Запуск:

  ```bash
  go test ./...
  ```

