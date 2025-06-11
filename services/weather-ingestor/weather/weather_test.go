package weather_test

import (
	"testing"
	"time"
	"weather-ingestor/weather"

	"github.com/google/uuid"
)

func TestCreateWeather(t *testing.T) {
	agriUnitID := uuid.New()
	input := weather.WeatherValue{
		Latitude:           48.8566,
		Longitude:          2.3522,
		Temperature:        18.5,
		Humidity:           60,
		WindSpeed:          5.2,
		Clouds:             75,
		WeatherMain:        "Clouds",
		WeatherDesc:        "broken clouds",
		AgriculturalUnitId: agriUnitID,
	}

	startTime := time.Now()
	w := weather.CreateWeather(input)

	if w.ID == uuid.Nil {
		t.Errorf("Expected a valid UUID, got Nil UUID")
	}

	if w.CreatedAt.Before(startTime) || time.Since(w.CreatedAt) > time.Second {
		t.Errorf("Unexpected CreatedAt timestamp: %v", w.CreatedAt)
	}

	if !w.UpdatedAt.Equal(w.CreatedAt) {
		t.Errorf("Expected UpdatedAt to equal CreatedAt")
	}

	if w.ArchivedAt != nil {
		t.Errorf("Expected ArchivedAt to be nil, got %v", w.ArchivedAt)
	}

	if w.Latitude != input.Latitude {
		t.Errorf("Expected Latitude %f, got %f", input.Latitude, w.Latitude)
	}
	if w.Longitude != input.Longitude {
		t.Errorf("Expected Longitude %f, got %f", input.Longitude, w.Longitude)
	}
	if w.Temperature != input.Temperature {
		t.Errorf("Expected Temperature %f, got %f", input.Temperature, w.Temperature)
	}
	if w.Humidity != input.Humidity {
		t.Errorf("Expected Humidity %d, got %d", input.Humidity, w.Humidity)
	}
	if w.WindSpeed != input.WindSpeed {
		t.Errorf("Expected WindSpeed %f, got %f", input.WindSpeed, w.WindSpeed)
	}
	if w.Clouds != input.Clouds {
		t.Errorf("Expected Clouds %d, got %d", input.Clouds, w.Clouds)
	}
	if w.WeatherMain != input.WeatherMain {
		t.Errorf("Expected WeatherMain '%s', got '%s'", input.WeatherMain, w.WeatherMain)
	}
	if w.WeatherDesc != input.WeatherDesc {
		t.Errorf("Expected WeatherDesc '%s', got '%s'", input.WeatherDesc, w.WeatherDesc)
	}
	if w.AgriculturalUnitId != input.AgriculturalUnitId {
		t.Errorf("Expected AgriculturalUnitId %s, got %s", input.AgriculturalUnitId, w.AgriculturalUnitId)
	}
}
