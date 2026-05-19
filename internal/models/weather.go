package models

import "encoding/json"

type OpenMeteoResponse struct {
	Latitude              float64      `json:"latitude"`
	Longitude             float64      `json:"longitude"`
	Generationtime_ms     float64      `json:"generationtime_ms"`
	Utc_offset_seconds    int          `json:"utc_offset_seconds"`
	Timezone              string       `json:"timezone"`
	Timezone_abbreviation string       `json:"timezone_abbreviation"`
	Elevation             float64      `json:"elevation"`
	Current_weather_units WeatherUnits `json:"current_weather_units"`
	Current_weather       Weather      `json:"current_weather"`
}

type WeatherUnits struct {
	Time          string `json:"time"`
	Interval      string `json:"interval"`
	Temperature   string `json:"temperature"`
	Windspeed     string `json:"windspeed"`
	Winddirection string `json:"winddirection"`
	Is_day        string `json:"is_day"`
	Weathercode   string `json:"weathercode"`
}

type Weather struct {
	Time          json.RawMessage `json:"time"`
	Interval      int             `json:"interval"`
	Temperature   float64         `json:"temperature"`
	Windspeed     float64         `json:"windspeed"`
	Winddirection float64         `json:"winddirection"`
	Is_day        float64         `json:"is_day"`
	Weathercode   float64         `json:"weathercode"`
}

type WeatherResponse struct {
	Temperature float64 `json:"temperature_c"`
	Time        string  `json:"time"`
	Source      string  `json:"source"` // "cache" или "api"
}
