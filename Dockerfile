# Используем многоэтапную сборку для уменьшения размера образа
FROM golang:1.24 AS builder

WORKDIR /app

# Копируем файлы зависимостей и загружаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект и собираем бинарный файл
COPY . .
RUN go build -o /app/main ./cmd

# Используем минимальный базовый образ для финального этапа
FROM alpine:latest

WORKDIR /app

# Копируем собранный бинарный файл из предыдущего этапа
COPY --from=builder /app/main .

# Отладка: вывод содержимого директории /app
RUN ls -la /app

# Команда запуска приложения
CMD ["/app/main"]
