package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

const googleMapsGeocodeURL = "https://maps.googleapis.com/maps/api/geocode/json"

var googleMapsHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

type googleMapsResult struct {
	Geometry struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	} `json:"geometry"`
	FormattedAddress string `json:"formatted_address"`
}

type googleMapsResponse struct {
	Results []googleMapsResult `json:"results"`
	Status  string             `json:"status"`
}

// GetCoordinatesFromGoogleMaps menggunakan Google Maps Geocoding API untuk mendapatkan koordinat
func GetCoordinatesFromGoogleMaps(location string) (lat float64, lon float64, err error) {
	if location == "" {
		return 0, 0, fmt.Errorf("location is required")
	}

	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return 0, 0, fmt.Errorf("GOOGLE_MAPS_API_KEY environment variable is not set")
	}

	u, err := url.Parse(googleMapsGeocodeURL)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse Google Maps URL: %w", err)
	}

	q := u.Query()
	q.Set("address", location)
	q.Set("key", apiKey)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create Google Maps request: %w", err)
	}

	resp, err := googleMapsHTTPClient.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to call Google Maps Geocoding API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("Google Maps Geocoding API returned status %d", resp.StatusCode)
	}

	var data googleMapsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, fmt.Errorf("failed to decode Google Maps response: %w", err)
	}

	if data.Status != "OK" {
		if data.Status == "ZERO_RESULTS" {
			return 0, 0, fmt.Errorf("location not found: %s", location)
		}
		return 0, 0, fmt.Errorf("Google Maps API error: %s", data.Status)
	}

	if len(data.Results) == 0 {
		return 0, 0, fmt.Errorf("location not found: %s", location)
	}

	result := data.Results[0]
	lat = result.Geometry.Location.Lat
	lon = result.Geometry.Location.Lng

	if lat == 0 && lon == 0 {
		return 0, 0, fmt.Errorf("invalid coordinates returned from Google Maps")
	}

	return lat, lon, nil
}
