package weather

import (
	"database/sql"
	"regexp"
	"testing"
	"time"
	"weather-ingestor/misc"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestWeatherToSqlViewAndBack(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)
	archived := now.Add(-24 * time.Hour)

	weather := Weather{
		ID:                 uuid.New(),
		CreatedAt:          now,
		UpdatedAt:          now,
		ArchivedAt:         &archived,
		Latitude:           48.8566,
		Longitude:          2.3522,
		Temperature:        23.5,
		Humidity:           65,
		WindSpeed:          5.2,
		Clouds:             75,
		WeatherMain:        "Clouds",
		WeatherDesc:        "scattered clouds",
		AgriculturalUnitId: uuid.New(),
	}

	sqlView := WeatherToSqlView(weather)

	if weather.ID != sqlView.ID {
		t.Errorf("ID mismatch: got %v want %v", sqlView.ID, weather.ID)
	}
	if !weather.CreatedAt.Equal(sqlView.CreatedAt) {
		t.Errorf("CreatedAt mismatch: got %v want %v", sqlView.CreatedAt, weather.CreatedAt)
	}
	if !weather.UpdatedAt.Equal(sqlView.UpdatedAt) {
		t.Errorf("UpdatedAt mismatch: got %v want %v", sqlView.UpdatedAt, weather.UpdatedAt)
	}
	if !sqlView.ArchivedAt.Valid {
		t.Error("ArchivedAt should be valid but is not")
	} else if !weather.ArchivedAt.Equal(sqlView.ArchivedAt.Time) {
		t.Errorf("ArchivedAt mismatch: got %v want %v", sqlView.ArchivedAt.Time, *weather.ArchivedAt)
	}
	if weather.Latitude != sqlView.Latitude {
		t.Errorf("Latitude mismatch: got %v want %v", sqlView.Latitude, weather.Latitude)
	}
	if weather.Longitude != sqlView.Longitude {
		t.Errorf("Longitude mismatch: got %v want %v", sqlView.Longitude, weather.Longitude)
	}
	if weather.Temperature != sqlView.Temperature {
		t.Errorf("Temperature mismatch: got %v want %v", sqlView.Temperature, weather.Temperature)
	}
	if weather.Humidity != sqlView.Humidity {
		t.Errorf("Humidity mismatch: got %v want %v", sqlView.Humidity, weather.Humidity)
	}
	if weather.WindSpeed != sqlView.WindSpeed {
		t.Errorf("WindSpeed mismatch: got %v want %v", sqlView.WindSpeed, weather.WindSpeed)
	}
	if weather.Clouds != sqlView.Clouds {
		t.Errorf("Clouds mismatch: got %v want %v", sqlView.Clouds, weather.Clouds)
	}
	if weather.WeatherMain != sqlView.WeatherMain {
		t.Errorf("WeatherMain mismatch: got %v want %v", sqlView.WeatherMain, weather.WeatherMain)
	}
	if weather.WeatherDesc != sqlView.WeatherDesc {
		t.Errorf("WeatherDesc mismatch: got %v want %v", sqlView.WeatherDesc, weather.WeatherDesc)
	}
	if weather.AgriculturalUnitId != sqlView.AgriculturalUnitId {
		t.Errorf("AgriculturalUnitId mismatch: got %v want %v", sqlView.AgriculturalUnitId, weather.AgriculturalUnitId)
	}

	converted, err := WeatherFromSqlView(sqlView)
	if err != nil {
		t.Fatalf("WeatherFromSqlView returned error: %v", err)
	}

	if weather.ID != converted.ID {
		t.Errorf("Converted ID mismatch: got %v want %v", converted.ID, weather.ID)
	}
	if !weather.CreatedAt.Equal(converted.CreatedAt) {
		t.Errorf("Converted CreatedAt mismatch: got %v want %v", converted.CreatedAt, weather.CreatedAt)
	}
	if !weather.UpdatedAt.Equal(converted.UpdatedAt) {
		t.Errorf("Converted UpdatedAt mismatch: got %v want %v", converted.UpdatedAt, weather.UpdatedAt)
	}
	if converted.ArchivedAt == nil {
		t.Error("Converted ArchivedAt is nil, want non-nil")
	} else if !weather.ArchivedAt.Equal(*converted.ArchivedAt) {
		t.Errorf("Converted ArchivedAt mismatch: got %v want %v", *converted.ArchivedAt, *weather.ArchivedAt)
	}
	if weather.Latitude != converted.Latitude {
		t.Errorf("Converted Latitude mismatch: got %v want %v", converted.Latitude, weather.Latitude)
	}
	if weather.Longitude != converted.Longitude {
		t.Errorf("Converted Longitude mismatch: got %v want %v", converted.Longitude, weather.Longitude)
	}
	if weather.Temperature != converted.Temperature {
		t.Errorf("Converted Temperature mismatch: got %v want %v", converted.Temperature, weather.Temperature)
	}
	if weather.Humidity != converted.Humidity {
		t.Errorf("Converted Humidity mismatch: got %v want %v", converted.Humidity, weather.Humidity)
	}
	if weather.WindSpeed != converted.WindSpeed {
		t.Errorf("Converted WindSpeed mismatch: got %v want %v", converted.WindSpeed, weather.WindSpeed)
	}
	if weather.Clouds != converted.Clouds {
		t.Errorf("Converted Clouds mismatch: got %v want %v", converted.Clouds, weather.Clouds)
	}
	if weather.WeatherMain != converted.WeatherMain {
		t.Errorf("Converted WeatherMain mismatch: got %v want %v", converted.WeatherMain, weather.WeatherMain)
	}
	if weather.WeatherDesc != converted.WeatherDesc {
		t.Errorf("Converted WeatherDesc mismatch: got %v want %v", converted.WeatherDesc, weather.WeatherDesc)
	}
	if weather.AgriculturalUnitId != converted.AgriculturalUnitId {
		t.Errorf("Converted AgriculturalUnitId mismatch: got %v want %v", converted.AgriculturalUnitId, weather.AgriculturalUnitId)
	}
}

func TestWeatherToSqlViewArchivedNil(t *testing.T) {
	now := time.Now()

	weather := Weather{
		ID:                 uuid.New(),
		CreatedAt:          now,
		UpdatedAt:          now,
		ArchivedAt:         nil,
		Latitude:           0,
		Longitude:          0,
		Temperature:        0,
		Humidity:           0,
		WindSpeed:          0,
		Clouds:             0,
		WeatherMain:        "",
		WeatherDesc:        "",
		AgriculturalUnitId: uuid.New(),
	}

	sqlView := WeatherToSqlView(weather)

	if sqlView.ArchivedAt.Valid {
		t.Error("ArchivedAt.Valid should be false when ArchivedAt is nil")
	}

	converted, err := WeatherFromSqlView(sqlView)
	if err != nil {
		t.Fatalf("WeatherFromSqlView returned error: %v", err)
	}

	if converted.ArchivedAt != nil {
		t.Errorf("Converted ArchivedAt should be nil, got %v", *converted.ArchivedAt)
	}
}

func TestInsertOrUpdate_Weather_Success(t *testing.T) {
	mockQuerierInstance, sqlMock, err := misc.NewMockQuerier(t)
	if err != nil {
		t.Fatalf("failed to create mock querier: %v", err)
	}
	defer mockQuerierInstance.Db.Close()

	storage := NewWeatherStorage(mockQuerierInstance)

	id := uuid.New()
	agriID := uuid.New()
	now := time.Now().Truncate(time.Millisecond)

	w := Weather{
		ID:                 id,
		CreatedAt:          now.Add(-24 * time.Hour),
		UpdatedAt:          now,
		ArchivedAt:         nil,
		Latitude:           50.0,
		Longitude:          3.0,
		Temperature:        21.5,
		Humidity:           75,
		WindSpeed:          4.5,
		Clouds:             20,
		WeatherMain:        "Rain",
		WeatherDesc:        "light rain",
		AgriculturalUnitId: agriID,
	}

	expectedSQL := "INSERT INTO weather (id,created_at,updated_at,archived_at,latitude,longitude,temperature,humidity,wind_speed,clouds,weather_main,weather_desc,agricultural_unit_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) ON CONFLICT (id) DO UPDATE SET updated_at = EXCLUDED.updated_at, archived_at = EXCLUDED.archived_at, latitude = EXCLUDED.latitude, longitude = EXCLUDED.longitude, temperature = EXCLUDED.temperature, humidity = EXCLUDED.humidity, wind_speed = EXCLUDED.wind_speed, clouds = EXCLUDED.clouds, weather_main = EXCLUDED.weather_main, weather_desc = EXCLUDED.weather_desc, agricultural_unit_id = EXCLUDED.agricultural_unit_id"

	sqlMock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
		WithArgs(
			w.ID,
			w.CreatedAt,
			w.UpdatedAt,
			sql.NullTime{Time: time.Time{}, Valid: false},
			w.Latitude,
			w.Longitude,
			w.Temperature,
			w.Humidity,
			w.WindSpeed,
			w.Clouds,
			w.WeatherMain,
			w.WeatherDesc,
			w.AgriculturalUnitId,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.InsertOrUpdate(w)
	if err != nil {
		t.Fatalf("InsertOrUpdate returned unexpected error: %v", err)
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
