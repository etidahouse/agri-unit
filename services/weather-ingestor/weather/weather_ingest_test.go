package weather_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"weather-ingestor/weather"

	"github.com/google/uuid"
)

type MockAgriUnitStorage struct{}

func (m *MockAgriUnitStorage) SelectAll() ([]weather.AgriculturalUnit, error) {
	return []weather.AgriculturalUnit{
		{
			ID:        uuid.New(),
			Latitude:  45.76,
			Longitude: 4.85,
		},
	}, nil
}

type MockWeatherStorage struct {
	Called bool
	Last   weather.Weather
}

func (m *MockWeatherStorage) InsertOrUpdate(w weather.Weather) error {
	m.Called = true
	m.Last = w
	return nil
}

func TestWeatherFetcher_Run_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"main": map[string]interface{}{
				"temp":     22.3,
				"humidity": 55,
			},
			"wind": map[string]interface{}{
				"speed": 4.5,
			},
			"clouds": map[string]interface{}{
				"all": 30,
			},
			"weather": []map[string]interface{}{
				{"main": "Cloudy", "description": "partly cloudy"},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	mockWeatherStorage := &MockWeatherStorage{}
	mockAgriUnitStorage := &MockAgriUnitStorage{}

	fetcher := weather.NewWeatherFetcher(server.URL, "dummy", mockWeatherStorage, mockAgriUnitStorage)

	err := fetcher.HandleWeatherIngest()
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if !mockWeatherStorage.Called {
		t.Errorf("InsertOrUpdate was not called")
	}

	got := mockWeatherStorage.Last
	if got.WeatherMain != "Cloudy" {
		t.Errorf("expected WeatherMain 'Cloudy', got '%s'", got.WeatherMain)
	}
	if got.WeatherDesc != "partly cloudy" {
		t.Errorf("expected WeatherDesc 'partly cloudy', got '%s'", got.WeatherDesc)
	}
	if got.Temperature != 22.3 {
		t.Errorf("expected Temperature 22.3, got %f", got.Temperature)
	}
	if got.Humidity != 55 {
		t.Errorf("expected Humidity 55, got %d", got.Humidity)
	}
}
