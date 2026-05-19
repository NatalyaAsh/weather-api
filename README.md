# Weather API

Микросервис на Go для получения текущей погоды с кешированием, метриками и Docker-контейнеризацией.

## 📋 Описание

Сервис получает данные о погоде из публичного API [Open-Meteo](https://open-meteo.com/), кеширует их на заданное время и отдаёт через собственный HTTP API.

### Особенности

- ✅ Чистая архитектура (слои: api, cache, client, config, models)
- ✅ In-memory кеш с TTL
- ✅ Конфигурация через YAML
- ✅ Graceful shutdown
- ✅ JSON-логирование (slog)
- ✅ Prometheus метрики
- ✅ Юнит-тесты
- ✅ Docker + docker-compose
- ✅ CI (GitHub Actions)

## 🚀 Быстрый старт

### Локальный запуск

```bash
# Склонировать репозиторий
git clone https://github.com/your-username/weather-api.git
cd weather-api

# Установить зависимости
go mod download

# Запустить сервер
go run cmd/main.go
```
### Проверка работы
```bash
# Получить погоду
curl http://localhost:8080/weather

# Проверить здоровье сервиса
curl http://localhost:8080/health

# Посмотреть метрики
curl http://localhost:8080/metrics
```
## Docker
```bash
# Собрать и запустить
docker-compose up --build

# Остановить
docker-compose down
```
## 📁 Структура проекта
```text
weather-api/
├── cmd/
│   └── main.go                 # Точка входа
├── internal/
│   ├── api/                    # HTTP-обработчики
│   ├── cache/                  # Кеш с TTL
│   ├── client/                 # Клиент для Open-Meteo
│   ├── config/                 # Загрузка конфига
│   └── models/                 # Структуры данных
├── .github/workflows/          # CI/CD
├── config.yml                  # Конфигурация
├── Dockerfile
├── docker-compose.yml
└── go.mod
```
## ⚙️ Конфигурация
Файл config.yml:
```yaml
port: "8080"                      # Порт сервера
cache_ttl_seconds: 300            # Время жизни кеша (секунды)
city_lat: 55.03                   # Широта (Новосибирск)
city_lon: 82.92                   # Долгота
city_name: "Novosibirsk"          # Название города
open_meteo_url: "https://api.open-meteo.com/v1/forecast"
```
## 📊 API
`GET /weather` - Возвращает текущую погоду.

Пример ответа:
```json
{
  "temperature_c": 11.4,
  "time": "2026-05-19 15:30:00",
  "source": "api"
}
```
Поле source показывает, откуда взяты данные: api (запрос к Open-Meteo) или cache (из кеша).

`GET /health` - Проверка работоспособности.

Пример ответа:
```json
{
  "status": "ok"
}
```
`GET /metrics` - Prometheus метрики (количество запросов, длительность и т.д.).

## 🧪 Тесты
```bash
go test ./...
📈 Метрики
Имя метрики	Тип	Описание
weather_requests_total	Counter	Количество запросов погоды (с тегом source)
weather_request_duration_seconds	Histogram	Длительность запросов
```
