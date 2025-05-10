# Kvant Task API

REST API для управления пользователями и их заказами на Go + PostgreSQL.

![Go](https://img.shields.io/badge/Go-1.24-blue) ![Gin](https://img.shields.io/badge/Gin-1.10-green) ![GORM](https://img.shields.io/badge/GORM-1.26-orange) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue)

---

## 🚀 Быстрый старт

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/Solarmatte/rest-api-users-orders-go.git
   cd rest-api-users-orders-go
   ```

2. Создайте образец окружения .env:

   ```bash
   cp .env.example .env
   ```

3. Запустите приложение с помощью Docker Compose:
   ```bash
   docker-compose up --build
   ```

4. Приложение будет доступно по адресу: `http://localhost:8080`

---

## 📚 Документация API

Документация доступна по адресу: `http://localhost:8080/swagger/index.html`

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

---

## 🧪 Тестирование

Для запуска тестов выполните:
```bash
   go test ./...
```

---

## 🛠️ Базовые команды

### Сборка приложения
```bash
   go build -o main ./cmd
```

### Запуск приложения локально
```bash
   go run ./cmd
```

### Запуск через Docker Compose
```bash
   docker-compose up --build
```

### Остановка через Docker Compose
```bash
   docker-compose down
```

### Обновление Swagger-документации
```bash
   swag init -g cmd/main.go -o docs --parseDependency --parseInternal --parseDepth 3
```

