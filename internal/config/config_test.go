package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Создаём временный YAML файл
	content := `port: "9090"
cache_ttl_seconds: 100
city_lat: 55.75
city_lon: 37.62
city_name: "Moscow"
open_meteo_url: "https://api.open-meteo.com/v1/forecast"`

	tmpfile, err := os.CreateTemp("", "config*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Загружаем конфиг
	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Проверяем значения
	if cfg.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", cfg.Port)
	}

	if cfg.CacheTTLSeconds != 100 {
		t.Errorf("Expected CacheTTLSeconds 100, got %d", cfg.CacheTTLSeconds)
	}

	if cfg.CacheTTL() != 100*time.Second {
		t.Errorf("Expected CacheTTL 100s, got %v", cfg.CacheTTL())
	}

	if cfg.CityLat != 55.75 {
		t.Errorf("Expected CityLat 55.75, got %f", cfg.CityLat)
	}

	if cfg.CityLon != 37.62 {
		t.Errorf("Expected CityLon 37.62, got %f", cfg.CityLon)
	}

	if cfg.CityName != "Moscow" {
		t.Errorf("Expected CityName 'Moscow', got %s", cfg.CityName)
	}

	if cfg.OpenMeteoURL != "https://api.open-meteo.com/v1/forecast" {
		t.Errorf("Expected OpenMeteoURL '%s', got '%s'",
			"https://api.open-meteo.com/v1/forecast", cfg.OpenMeteoURL)
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/file.yml")
	if err == nil {
		t.Error("Expected error when file not found, got nil")
	}
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "config*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Некорректный YAML
	tmpfile.WriteString("port: [8080\n")
	tmpfile.Close()

	_, err = Load(tmpfile.Name())
	if err == nil {
		t.Error("Expected error with invalid YAML, got nil")
	}
}
