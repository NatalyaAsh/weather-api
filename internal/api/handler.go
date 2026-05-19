package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"weather-api/internal/cache"
	"weather-api/internal/client"
	"weather-api/internal/config"
	"weather-api/internal/models"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "weather_requests_total",
			Help: "Total number of weather requests",
		},
		[]string{"source"}, // "cache" или "api"
	)
)

type Handler struct {
	Config *config.Config
	Cache  *cache.Cache
	Client *client.WeatherClient
}

func NewHandler(cfg *config.Config, cache *cache.Cache, weatherClient *client.WeatherClient) *Handler {
	return &Handler{
		Config: cfg,
		Cache:  cache,
		Client: weatherClient,
	}
}

func (h *Handler) GetWeather(w http.ResponseWriter, r *http.Request) {
	key := "novosibirsk"

	slog.Info("Weather request received")

	// Проверяем кеш
	if cached, found := h.Cache.Get(key); found {
		slog.Info("Returning from cache")
		requestsTotal.WithLabelValues("cache").Inc()
		response := cached.(models.WeatherResponse)
		response.Source = "cache"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Нет в кеше — идём в Open-Meteo через клиент
	slog.Info("Fetching from Open-Meteo API")
	requestsTotal.WithLabelValues("api").Inc()
	weatherData, err := h.Client.FetchWeather(h.Config.CityLat, h.Config.CityLon)
	if err != nil {
		slog.Error("Failed to fetch weather", "error", err)
		http.Error(w, "Failed to fetch weather: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Парсим время (используем json.RawMessage)
	var timeStr string
	err = json.Unmarshal(weatherData.Current_weather.Time, &timeStr)
	if err != nil {
		http.Error(w, "Failed to parse time: "+err.Error(), http.StatusInternalServerError)
		return
	}

	parsedTime, err := time.Parse("2006-01-02T15:04", timeStr)
	if err != nil {
		http.Error(w, "Failed to parse time string: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	response := models.WeatherResponse{
		Temperature: weatherData.Current_weather.Temperature,
		Time:        parsedTime.Format("2006-01-02 15:04:05"),
		Source:      "api",
	}

	// Сохраняем в кеш
	h.Cache.Set(key, response)

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	slog.Info("Weather fetched successfully", "temperature", weatherData.Current_weather.Temperature)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
