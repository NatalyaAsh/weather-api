package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseOpenMeteoTime(t *testing.T) {
	// JSON с нестандартным форматом времени
	jsonData := `{
        "time": "2026-05-19T15:30",
        "interval": 900,
        "temperature": 11.4,
        "windspeed": 5.4,
        "winddirection": 340,
        "is_day": 1,
        "weathercode": 0
    }`

	var weather Weather
	err := json.Unmarshal([]byte(jsonData), &weather)
	if err != nil {
		t.Fatal("Failed to unmarshal JSON:", err)
	}

	// Парсим время из RawMessage
	var timeStr string
	err = json.Unmarshal(weather.Time, &timeStr)
	if err != nil {
		t.Fatal("Failed to parse time field:", err)
	}

	expectedFormat := "2026-05-19T15:30"
	if timeStr != expectedFormat {
		t.Errorf("Expected time '%s', got '%s'", expectedFormat, timeStr)
	}

	// Парсим в time.Time
	parsedTime, err := time.Parse("2006-01-02T15:04", timeStr)
	if err != nil {
		t.Fatal("Failed to parse time string:", err)
	}

	// Проверяем компоненты времени
	if parsedTime.Year() != 2026 {
		t.Errorf("Expected year 2026, got %d", parsedTime.Year())
	}
	if parsedTime.Month() != time.May {
		t.Errorf("Expected month May, got %s", parsedTime.Month())
	}
	if parsedTime.Day() != 19 {
		t.Errorf("Expected day 19, got %d", parsedTime.Day())
	}
	if parsedTime.Hour() != 15 {
		t.Errorf("Expected hour 15, got %d", parsedTime.Hour())
	}
	if parsedTime.Minute() != 30 {
		t.Errorf("Expected minute 30, got %d", parsedTime.Minute())
	}
}

func TestWeatherResponseJSON(t *testing.T) {
	response := WeatherResponse{
		Temperature: 11.4,
		Time:        "2026-05-19 15:30:00",
		Source:      "api",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatal("Failed to marshal WeatherResponse:", err)
	}

	var parsed WeatherResponse
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		t.Fatal("Failed to unmarshal WeatherResponse:", err)
	}

	if parsed.Temperature != 11.4 {
		t.Errorf("Expected temperature 11.4, got %f", parsed.Temperature)
	}
	if parsed.Source != "api" {
		t.Errorf("Expected source 'api', got '%s'", parsed.Source)
	}
}
