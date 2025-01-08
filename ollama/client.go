package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"learn-go/weather"
)

type Client struct {
	httpClient *http.Client
	weatherSvc weather.Service
	baseURL    string
}

func NewClient(weatherSvc weather.Service) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		weatherSvc: weatherSvc,
		baseURL:    "http://localhost:11434/api",
	}
}

type OllamaRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
	Options  *Options      `json:"options,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type WeatherResponse struct {
	Weather     *weather.WeatherData `json:"weather"`
	Description string               `json:"description"`
}

type Options struct {
	Temperature float64 `json:"temperature,omitempty"`
	Seed        int     `json:"seed,omitempty"`
}

type OllamaResponse struct {
	Model         string      `json:"model"`
	CreatedAt     string      `json:"created_at"`
	Message       ChatMessage `json:"message"`
	Done          bool        `json:"done"`
	TotalDuration int64       `json:"total_duration"`
}

// Add method to extract location from natural query
func (c *Client) ExtractLocation(query string) (string, error) {
	query = strings.TrimSpace(query)

	req := OllamaRequest{
		Model:  "phi4",
		Stream: false,
		Options: &Options{
			Temperature: 0.1,
			Seed:        42,
		},
		Messages: []ChatMessage{
			{
				Role: "system",
				Content: `You are a location extractor. Extract only the city and state/country from the query. 
                         Format: "City, State" or "City, Country". If no location is found, say "no location".`,
			},
			{
				Role:    "user",
				Content: query,
			},
		},
	}

	location, err := c.getAIResponse(req, 3)
	if err != nil {
		return "", fmt.Errorf("extracting location: %w", err)
	}

	if location == "no location" {
		return "", fmt.Errorf("no location found in query")
	}

	return location, nil
}

// Update GetWeather to handle natural language
func (c *Client) GetWeather(query string, maxRetries int) (*WeatherResponse, error) {
	// First extract location from query
	location, err := c.ExtractLocation(query)
	if err != nil {
		return nil, err
	}

	// Then get weather data using existing method
	return c.GetWeatherData(location, maxRetries)
}

func (c *Client) GetWeatherData(location string, maxRetries int) (*WeatherResponse, error) {
	weatherData, err := c.weatherSvc.GetWeather(location)
	if err != nil {
		return nil, err
	}

	weatherJSON, err := json.Marshal(weatherData)
	if err != nil {
		return nil, fmt.Errorf("marshaling weather data: %w", err)
	}

	req := OllamaRequest{
		Model:  "phi4", // phi4 is the default model
		Stream: false,
		Options: &Options{
			Temperature: 0.7,
			Seed:        42,
		},
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: "You are a weather assistant. Provide a natural, concise description of the weather conditions. Focus on temperature, conditions, and humidity.",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Describe the weather in %s based on this data: %s", location, string(weatherJSON)),
			},
		},
	}

	description, err := c.getAIResponse(req, maxRetries)
	if err != nil {
		return nil, fmt.Errorf("getting AI response: %w", err)
	}

	return &WeatherResponse{
		Weather:     weatherData,
		Description: description,
	}, nil
}

func (c *Client) getAIResponse(req OllamaRequest, maxRetries int) (string, error) {
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Second * time.Duration(1<<uint(attempt)))
		}

		jsonData, err := json.Marshal(req)
		if err != nil {
			return "", fmt.Errorf("marshaling request: %w", err)
		}

		request, err := http.NewRequest("POST", c.baseURL+"/chat", bytes.NewBuffer(jsonData))
		if err != nil {
			return "", fmt.Errorf("creating request: %w", err)
		}

		request.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(request)
		if err != nil {
			if attempt == maxRetries {
				return "", fmt.Errorf("failed to connect to ollama server (is it running?): %w", err)
			}
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return "", fmt.Errorf("ollama server error (status %d): %s", resp.StatusCode, string(body))
		}

		var aiResp OllamaResponse
		if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
			if attempt == maxRetries {
				return "", fmt.Errorf("decoding response: %w", err)
			}
			continue
		}

		if !aiResp.Done {
			continue
		}

		if aiResp.Message.Content == "" {
			if attempt == maxRetries {
				return "", fmt.Errorf("empty response from model")
			}
			continue
		}

		return aiResp.Message.Content, nil
	}

	return "", fmt.Errorf("max retries reached - please ensure ollama is running and the phi model is installed")
}
