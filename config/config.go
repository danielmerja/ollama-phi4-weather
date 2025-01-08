package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	APIKey         string
	OllamaURL      string
	OllamaModel    string
	RateLimit      float64
	AllowedOrigins []string
}

func Load() (*Config, error) {
	rateLimit, _ := strconv.ParseFloat(getEnvOrDefault("RATE_LIMIT", "1"), 64)

	return &Config{
		Port:           getEnvOrDefault("PORT", "8080"),
		APIKey:         os.Getenv("API_KEY"),
		OllamaURL:      getEnvOrDefault("OLLAMA_URL", "http://localhost:11434/api"),
		OllamaModel:    getEnvOrDefault("OLLAMA_MODEL", "phi4"),
		RateLimit:      rateLimit,
		AllowedOrigins: []string{"*"}, // Configure as needed
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
