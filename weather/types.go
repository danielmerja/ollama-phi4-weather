package weather

// WeatherData represents the processed weather information
type WeatherData struct {
	Temperature         float64 `json:"temperature"`
	FeelsLike           float64 `json:"feels_like"`
	Conditions          string  `json:"conditions"`
	Humidity            int     `json:"humidity"`
	WindSpeed           float64 `json:"wind_speed"`
	WindSpeedUnit       string  `json:"wind_speed_unit"`
	WindDirection       string  `json:"wind_direction"`
	WindGust            float64 `json:"wind_gust"`
	Visibility          float64 `json:"visibility"`
	Pressure            float64 `json:"pressure"`
	DewPoint            float64 `json:"dew_point"`
	UVIndex             float64 `json:"uv_index"`
	CloudCover          int     `json:"cloud_cover"`
	PrecipitationChance int     `json:"precipitation_chance"`
	QualityControl      string  `json:"quality_control"`
	Timestamp           string  `json:"timestamp"`
}

// API response structures
type PointsResponse struct {
	Properties struct {
		GridID              string           `json:"gridId"`
		GridX               int              `json:"gridX"`
		GridY               int              `json:"gridY"`
		RelativeLocation    RelativeLocation `json:"relativeLocation"`
		Forecast            string           `json:"forecast"`
		ForecastHourly      string           `json:"forecastHourly"`
		ObservationStations string           `json:"observationStations"`
	} `json:"properties"`
}

type RelativeLocation struct {
	Properties struct {
		City     string   `json:"city"`
		State    string   `json:"state"`
		Distance Distance `json:"distance"`
	} `json:"properties"`
}

type Distance struct {
	Value    float64 `json:"value"`
	UnitCode string  `json:"unitCode"`
}

type StationsResponse struct {
	Features []struct {
		Properties struct {
			StationIdentifier string `json:"stationIdentifier"`
			Name              string `json:"name"`
			TimeZone          string `json:"timeZone"`
			Status            string `json:"status"`
		} `json:"properties"`
	} `json:"features"`
}

type NWSResponse struct {
	Properties struct {
		Temperature struct {
			Value          float64 `json:"value"`
			UnitCode       string  `json:"unitCode"`
			QualityControl string  `json:"qualityControl"`
		} `json:"temperature"`
		RelativeHumidity struct {
			Value          float64 `json:"value"`
			UnitCode       string  `json:"unitCode"`
			QualityControl string  `json:"qualityControl"`
		} `json:"relativeHumidity"`
		WindSpeed struct {
			Value    float64 `json:"value"`
			UnitCode string  `json:"unitCode"`
		} `json:"windSpeed"`
		TextDescription string `json:"textDescription"`
		Timestamp       string `json:"timestamp"`
	} `json:"properties"`
}

// ... rest of the API types ...
