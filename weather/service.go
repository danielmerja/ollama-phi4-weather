package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Service interface {
	GetWeather(location string) (*WeatherData, error)
}

type NWSService struct {
	client    *http.Client
	userAgent string
}

func NewNWSService() *NWSService {
	return &NWSService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		userAgent: userAgent,
	}
}

func (s *NWSService) GetWeather(location string) (*WeatherData, error) {
	coords, err := s.validateLocation(location)
	if err != nil {
		return nil, err
	}

	points, err := s.getPoints(coords)
	if err != nil {
		return nil, fmt.Errorf("getting points data: %w", err)
	}

	station, err := s.findNearestStation(points)
	if err != nil {
		return nil, fmt.Errorf("finding station: %w", err)
	}

	return s.getWeatherData(station)
}

func (s *NWSService) validateLocation(location string) (*GeoLocation, error) {
	geo, err := geocode(location)
	if err != nil {
		return nil, fmt.Errorf("geocoding location: %w", err)
	}

	lat, _ := strconv.ParseFloat(geo.Lat, 64)
	lon, _ := strconv.ParseFloat(geo.Lon, 64)

	if !isUSLocation(lat, lon) {
		return nil, newLocationError(location)
	}

	return geo, nil
}

func (s *NWSService) getPoints(geo *GeoLocation) (*PointsResponse, error) {
	url := fmt.Sprintf("%s/points/%s,%s", nwsBaseURL, geo.Lat, geo.Lon)
	points := &PointsResponse{}

	if err := s.makeRequest(url, points); err != nil {
		return nil, err
	}

	return points, nil
}

func (s *NWSService) findNearestStation(points *PointsResponse) (string, error) {
	url := fmt.Sprintf("%s/gridpoints/%s/%d,%d/stations",
		nwsBaseURL,
		points.Properties.GridID,
		points.Properties.GridX,
		points.Properties.GridY)

	stations := &StationsResponse{}
	if err := s.makeRequest(url, stations); err != nil {
		return "", err
	}

	if len(stations.Features) == 0 {
		return "", fmt.Errorf("no weather stations found")
	}

	return stations.Features[0].Properties.StationIdentifier, nil
}

func (s *NWSService) getWeatherData(stationID string) (*WeatherData, error) {
	url := fmt.Sprintf("%s/stations/%s/observations/latest", nwsBaseURL, stationID)

	var nwsResp NWSResponse
	if err := s.makeRequest(url, &nwsResp); err != nil {
		return nil, fmt.Errorf("getting observations: %w", err)
	}

	return convertResponse(&nwsResp), nil
}

func convertResponse(resp *NWSResponse) *WeatherData {
	tempF := (resp.Properties.Temperature.Value * 9 / 5) + 32

	return &WeatherData{
		Temperature:    tempF,
		Conditions:     resp.Properties.TextDescription,
		Humidity:       int(resp.Properties.RelativeHumidity.Value),
		WindSpeed:      resp.Properties.WindSpeed.Value,
		WindSpeedUnit:  resp.Properties.WindSpeed.UnitCode,
		QualityControl: resp.Properties.Temperature.QualityControl,
		Timestamp:      resp.Properties.Timestamp,
	}
}

func (s *NWSService) makeRequest(url string, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Accept", "application/geo+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}
