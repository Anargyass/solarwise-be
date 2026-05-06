package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"solar-backend/internal/config"
)

const defaultNasaPowerURL = "https://power.larc.nasa.gov/api/temporal/climatology/point"

var nasaHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

type nasaResponse struct {
	Properties struct {
		Parameter struct {
			AllSkySfcSwDwn map[string]json.RawMessage `json:"ALLSKY_SFC_SW_DWN"`
			T2M            map[string]json.RawMessage `json:"T2M"`
		} `json:"parameter"`
	} `json:"properties"`
}

func GetClimatologyData(lat float64, lon float64) (radiasi float64, suhu float64, err error) {
	nasaURL := config.GetEnvOrDefault("NASA_POWER_URL", defaultNasaPowerURL)
	u, err := url.Parse(nasaURL)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse NASA URL: %w", err)
	}

	q := u.Query()
	q.Set("parameters", "ALLSKY_SFC_SW_DWN,T2M")
	q.Set("community", "RE")
	q.Set("longitude", strconv.FormatFloat(lon, 'f', 6, 64))
	q.Set("latitude", strconv.FormatFloat(lat, 'f', 6, 64))
	q.Set("format", "JSON")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create NASA request: %w", err)
	}

	resp, err := nasaHTTPClient.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to call NASA POWER: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("NASA POWER returned status %d", resp.StatusCode)
	}

	var data nasaResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, fmt.Errorf("failed to decode NASA POWER response: %w", err)
	}

	rawRadiasi, ok := data.Properties.Parameter.AllSkySfcSwDwn["ANN"]
	if !ok {
		return 0, 0, fmt.Errorf("missing ANN radiation data")
	}
	rawSuhu, ok := data.Properties.Parameter.T2M["ANN"]
	if !ok {
		return 0, 0, fmt.Errorf("missing ANN temperature data")
	}

	if err := json.Unmarshal(rawRadiasi, &radiasi); err != nil {
		return 0, 0, fmt.Errorf("invalid ANN radiation value: %w", err)
	}
	if err := json.Unmarshal(rawSuhu, &suhu); err != nil {
		return 0, 0, fmt.Errorf("invalid ANN temperature value: %w", err)
	}

	if radiasi == -999.0 || suhu == -999.0 {
		return 0, 0, fmt.Errorf("Data not available for this location")
	}

	return radiasi, suhu, nil
}
