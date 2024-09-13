# Используем официальный образ Go как базовый
FROM golang:1.22.5-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы с зависимостями
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Используем минимальный образ для запуска приложения
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем собранное приложение из предыдущего этапа
COPY --from=builder /app/main .

# Запускаем приложение
CMD ["./main"]