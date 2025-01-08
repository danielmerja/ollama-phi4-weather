package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// GeoLocation represents the latitude and longitude coordinates
type GeoLocation struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// geocode converts a location string to coordinates using OpenStreetMap's Nominatim service
func geocode(location string) (*GeoLocation, error) {
	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Add("q", location)
	params.Add("format", "json")
	params.Add("limit", "1")

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating geocode request: %w", err)
	}

	// Required by Nominatim's usage policy
	req.Header.Set("User-Agent", "WeatherApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("geocoding request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding service returned status: %d", resp.StatusCode)
	}

	var results []GeoLocation
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("decoding geocode response: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("location not found: %s", location)
	}

	return &results[0], nil
}
