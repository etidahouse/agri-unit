package weather

import (
	"time"

	"github.com/google/uuid"
)

type Weather struct {
	ID                 uuid.UUID  `json:"id"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	ArchivedAt         *time.Time `json:"archivedAt,omitempty"`
	Latitude           float64    `json:"latitude"`
	Longitude          float64    `json:"longitude"`
	Temperature        float64    `json:"temperature"`
	Humidity           int        `json:"humidity"`
	WindSpeed          float64    `json:"wind_speed"`
	Clouds             int        `json:"clouds"`
	WeatherMain        string     `json:"weather_main"`
	WeatherDesc        string     `json:"weather_desc"`
	AgriculturalUnitId uuid.UUID  `json:"agricultural_unit_id"`
}

type WeatherValue struct {
	Latitude           float64
	Longitude          float64
	Temperature        float64
	Humidity           int
	WindSpeed          float64
	Clouds             int
	WeatherMain        string
	WeatherDesc        string
	AgriculturalUnitId uuid.UUID
}

func CreateWeather(value WeatherValue) Weather {
	now := time.Now()
	return Weather{
		ID:                 uuid.New(),
		CreatedAt:          now,
		UpdatedAt:          now,
		Latitude:           value.Latitude,
		Longitude:          value.Longitude,
		Temperature:        value.Temperature,
		Humidity:           value.Humidity,
		WindSpeed:          value.WindSpeed,
		Clouds:             value.Clouds,
		WeatherMain:        value.WeatherMain,
		WeatherDesc:        value.WeatherDesc,
		AgriculturalUnitId: value.AgriculturalUnitId,
	}
}
