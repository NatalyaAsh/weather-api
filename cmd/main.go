package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"weather-api/internal/api"
	"weather-api/internal/cache"
	"weather-api/internal/client"
	"weather-api/internal/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Настраиваем логгер
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Загружаем конфиг
	cfg, err := config.Load("config.yml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Создаём кеш
	weatherCache := cache.New(cfg.CacheTTL())

	// Создаём клиент для Open-Meteo
	weatherClient := client.New(cfg.OpenMeteoURL)

	// Создаём обработчик (передаём все зависимости)
	handler := api.NewHandler(cfg, weatherCache, weatherClient)

	// Регистрируем маршруты
	mux := http.NewServeMux()
	mux.HandleFunc("/weather", handler.GetWeather)
	mux.HandleFunc("/health", handler.HealthCheck)
	mux.Handle("/metrics", promhttp.Handler())

	// Создаём сервер
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		fmt.Println("Сервер запущен на http://localhost:" + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")

	// Таймаут на завершение текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("Server exited gracefully")
}
