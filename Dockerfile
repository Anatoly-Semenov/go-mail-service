# Этап сборки
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Установка необходимых зависимостей
RUN apk add --no-cache git

# Копирование файлов проекта
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o mail-service ./cmd/app

# Финальный этап
FROM alpine:latest

WORKDIR /app

# Копирование бинарного файла из этапа сборки
COPY --from=builder /app/mail-service .

# Запуск приложения
CMD ["./mail-service"] 