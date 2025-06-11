package weather

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type WeatherFetcher struct {
	apiURL          string
	apiKey          string
	weatherStorage  WeatherStorage
	agriUnitStorage AgriUnitStorage
}

func NewWeatherFetcher(apiURL, apiKey string, ws WeatherStorage, aus AgriUnitStorage) *WeatherFetcher {
	return &WeatherFetcher{
		apiURL:          apiURL,
		apiKey:          apiKey,
		weatherStorage:  ws,
		agriUnitStorage: aus,
	}
}

func (wf *WeatherFetcher) HandleWeatherIngest() error {
	units, err := wf.agriUnitStorage.SelectAll()
	if err != nil {
		return fmt.Errorf("failed to fetch agri units: %w", err)
	}

	for _, unit := range units {
		weather, err := wf.fetchWeather(unit.Latitude, unit.Longitude, unit.ID)
		if err != nil {
			fmt.Printf("failed to fetch weather for unit %v: %v\n", unit.ID, err)
			continue
		}

		if err := wf.weatherStorage.InsertOrUpdate(weather); err != nil {
			fmt.Printf("failed to save weather for unit %v: %v\n", unit.ID, err)
		}
	}

	return nil
}

func (wf *WeatherFetcher) fetchWeather(lat, lon float64, agriUnitId uuid.UUID) (Weather, error) {
	url := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s&units=metric", wf.apiURL, lat, lon, wf.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return Weather{}, fmt.Errorf("api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Weather{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var data struct {
		Main struct {
			Temp     float64 `json:"temp"`
			Humidity int     `json:"humidity"`
		} `json:"main"`
		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Weather []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
		} `json:"weather"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return Weather{}, fmt.Errorf("decoding failed: %w", err)
	}

	main := ""
	desc := ""
	if len(data.Weather) > 0 {
		main = data.Weather[0].Main
		desc = data.Weather[0].Description
	}

	weather := CreateWeather(WeatherValue{
		Latitude:           lat,
		Longitude:          lon,
		Temperature:        data.Main.Temp,
		Humidity:           data.Main.Humidity,
		WindSpeed:          data.Wind.Speed,
		Clouds:             data.Clouds.All,
		WeatherMain:        main,
		WeatherDesc:        desc,
		AgriculturalUnitId: agriUnitId,
	})

	return weather, nil
}
