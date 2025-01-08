package weather

const (
	nwsBaseURL = "https://api.weather.gov"
	userAgent  = "WeatherApp/1.0"
)

// USBounds defines the geographical boundaries for US territories
type USBounds struct {
	MinLat, MaxLat float64
	MinLon, MaxLon float64
}

var (
	continentalUS = USBounds{
		MinLat: 24.396308,
		MaxLat: 49.384358,
		MinLon: -125.000000,
		MaxLon: -66.934570,
	}
	alaska = USBounds{
		MinLat: 51.214183,
		MaxLat: 71.365162,
		MinLon: -179.148909,
		MaxLon: -130.977806,
	}
	hawaii = USBounds{
		MinLat: 18.910361,
		MaxLat: 22.236428,
		MinLon: -160.236068,
		MaxLon: -154.808063,
	}
)

func isUSLocation(lat, lon float64) bool {
	return (lat >= continentalUS.MinLat && lat <= continentalUS.MaxLat && lon >= continentalUS.MinLon && lon <= continentalUS.MaxLon) || // Continental
		(lat >= alaska.MinLat && lat <= alaska.MaxLat && lon >= alaska.MinLon && lon <= alaska.MaxLon) || // Alaska
		(lat >= hawaii.MinLat && lat <= hawaii.MaxLat && lon >= hawaii.MinLon && lon <= hawaii.MaxLon) // Hawaii
}
