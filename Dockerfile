# Используем golang 1.24 как базовый образ
FROM golang:1.24

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект внутрь контейнера
COPY . .

# Собираем главный бинарь из ./cmd
RUN go build -o main ./cmd

# Команда запуска приложения
CMD ["./main"]
