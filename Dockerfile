# Этап сборки
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o /weather-api ./cmd/main.go

# Этап запуска
FROM alpine:latest

# Устанавливаем CA-сертификаты для HTTPS-запросов
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из этапа сборки
COPY --from=builder /weather-api .

# Копируем конфиг
COPY config.yml .

# Открываем порт
EXPOSE 8080

# Запускаем
CMD ["./weather-api"]