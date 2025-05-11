# ── builder stage ───────────────────────────────────────────────────────────────
FROM golang:1.24 AS builder
WORKDIR /app

# Статическая сборка без CGO (бинарник не будет зависеть от glibc)
ENV CGO_ENABLED=0
ENV GOOS=linux

# Копируем зависимости и скачиваем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код и собираем бинарник
COPY . .
RUN go build -ldflags="-s -w" -o main ./cmd

# ── final stage ────────────────────────────────────────────────────────────────
FROM scratch
WORKDIR /app

# Копируем собранный бинарник из builder
COPY --from=builder /app/main .

# Если нужно HTTPS/SSL, можно добавить сертификаты:
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Точка входа
ENTRYPOINT ["/app/main"]
