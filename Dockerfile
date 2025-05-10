# Используем golang 1.24 образ как базовый
FROM golang:1.24

# Рабочая директория в контейнере
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект внутрь
COPY . .

# Собираем главный бинарь из ./cmd
RUN go build -o main ./cmd

# Команда запуска приложения
CMD ["./main"]
