package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"weather-api/internal/models"
)

type WeatherClient struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *WeatherClient {
	return &WeatherClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *WeatherClient) FetchWeather(lat, lon float64) (*models.OpenMeteoResponse, error) {
	slog.Info("Calling Open-Meteo API", "lat", lat, "lon", lon)
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&current_weather=true",
		c.baseURL, lat, lon)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		slog.Error("HTTP request failed", "error", err)
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("HTTP request failed", "error", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var weatherData models.OpenMeteoResponse
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		slog.Error("HTTP request failed", "error", err)
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &weatherData, nil
}
