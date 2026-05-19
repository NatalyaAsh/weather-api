package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port            string  `yaml:"port"`
	CacheTTLSeconds int     `yaml:"cache_ttl_seconds"`
	CityLat         float64 `yaml:"city_lat"`
	CityLon         float64 `yaml:"city_lon"`
	CityName        string  `yaml:"city_name"`
	OpenMeteoURL    string  `yaml:"open_meteo_url"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// CacheTTL возвращает время жизни кеша как time.Duration
func (c *Config) CacheTTL() time.Duration {
	return time.Duration(c.CacheTTLSeconds) * time.Second
}
