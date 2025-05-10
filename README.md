
# API для управления пользователями и заказами

## Описание проекта

Данный проект представляет собой REST API, написанный на Go, для управления списком пользователей и их заказами с использованием базы данных PostgreSQL. В API реализованы функционал создания, обновления, удаления и получения пользователей и заказов, а также JWT-авторизация для защиты маршрутов. Проект также содержит Swagger-документацию для удобства ознакомления с методами API.

---

## Требования

- Go 1.24+
- Docker и Docker Compose (для удобного запуска)
- PostgreSQL (если запускать без Docker)

---

## Быстрый запуск с Docker Compose

1. Скопируйте файл `.env.example` (если есть) или создайте `.env` в корне проекта и заполните переменные окружения:

```env
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=qwerty
POSTGRES_DB=kvant_db

JWT_SECRET=qwerty

SERVER_ADDRESS=:8080
```

2. Запустите сервисы командой:

```bash
docker-compose up --build
```

3. После запуска API будет доступен на [http://localhost:8080](http://localhost:8080), а Swagger UI — на [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

4. Для остановки контейнеров используйте:

```bash
docker-compose down
```

---

## Запуск локально без Docker

1. Установите и запустите PostgreSQL.

2. Создайте базу данных, если её нет:

```sql
CREATE DATABASE kvant_db;
```

3. Создайте файл `.env` с параметрами подключения к БД и конфигурацией JWT, аналогично примеру выше.

4. Установите зависимости и запустите сервер:

```bash
go mod download
go run cmd/main.go
```

---

## API Описание

Все эндпоинты, кроме `/auth/login` и создания пользователя, защищены JWT-токеном, который необходимо передавать в заголовке `Authorization` со схемой `Bearer`.

### Пользователи

- **POST /users** — создать пользователя  
  Тело (JSON):  
  ```json
  {
    "name": "John Doe",
    "email": "john.doe@example.com",
    "age": 30,
    "password": "securepassword"
  }
  ```
- **GET /users** — получить список с пагинацией и фильтрацией по возрасту  
  Параметры: `page`, `limit`, `min_age`, `max_age`  

- **GET /users/{id}** — получить пользователя по ID  
- **PUT /users/{id}** — обновить данные пользователя  
- **DELETE /users/{id}** — удалить пользователя  

### Заказы

- **POST /users/{user_id}/orders** — создать заказ для пользователя  
  Тело (JSON):  
  ```json
  {
    "product": "Laptop",
    "quantity": 1,
    "price": 1200.50
  }
  ```
- **GET /users/{user_id}/orders** — получить список заказов пользователя  

### Авторизация

- **POST /auth/login** — аутентификация и получение JWT  
  Тело (JSON):  
  ```json
  {
    "email": "john.doe@example.com",
    "password": "securepassword"
  }
  ```
  Ответ содержит токен:
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```

---

## Тестирование

1. Все тесты находятся в каталоге `tests`.

2. Для запуска тестов выполните:

```bash
go test ./tests/...
```

---

## Логирование

Проект ведёт логирование основных операций: создание, обновление, удаление пользователей и заказов для удобства мониторинга.

---

## Миграции и модели БД

При запуске приложения автоматически выполняется миграция моделей `users` и `orders` через GORM.

Схема таблиц:

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  age INT NOT NULL,
  password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  product VARCHAR(255) NOT NULL,
  quantity INT NOT NULL,
  price DECIMAL(10, 2) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Swagger документация

Документация доступна по адресу:  
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)



---

###@ Спасибо за использование проекта!
 