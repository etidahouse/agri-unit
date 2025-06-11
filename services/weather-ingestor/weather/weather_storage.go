package weather

import (
	"database/sql"
	"fmt"
	"time"
	"weather-ingestor/storage"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type WeatherSqlView struct {
	ID                 uuid.UUID    `db:"id"`
	CreatedAt          time.Time    `db:"created_at"`
	UpdatedAt          time.Time    `db:"updated_at"`
	ArchivedAt         sql.NullTime `db:"archived_at"`
	Latitude           float64      `db:"latitude"`
	Longitude          float64      `db:"longitude"`
	Temperature        float64      `db:"temperature"`
	Humidity           int          `db:"humidity"`
	WindSpeed          float64      `db:"wind_speed"`
	Clouds             int          `db:"clouds"`
	WeatherMain        string       `db:"weather_main"`
	WeatherDesc        string       `db:"weather_desc"`
	AgriculturalUnitId uuid.UUID    `db:"agricultural_unit_id"`
}

func WeatherToSqlView(w Weather) WeatherSqlView {
	var archived sql.NullTime
	if w.ArchivedAt != nil {
		archived = sql.NullTime{Time: *w.ArchivedAt, Valid: true}
	} else {
		archived = sql.NullTime{Valid: false}
	}

	return WeatherSqlView{
		ID:                 w.ID,
		CreatedAt:          w.CreatedAt,
		UpdatedAt:          w.UpdatedAt,
		ArchivedAt:         archived,
		Latitude:           w.Latitude,
		Longitude:          w.Longitude,
		Temperature:        w.Temperature,
		Humidity:           w.Humidity,
		WindSpeed:          w.WindSpeed,
		Clouds:             w.Clouds,
		WeatherMain:        w.WeatherMain,
		WeatherDesc:        w.WeatherDesc,
		AgriculturalUnitId: w.AgriculturalUnitId,
	}
}

func WeatherFromSqlView(sqlView WeatherSqlView) (Weather, error) {

	var archivedAt *time.Time
	if sqlView.ArchivedAt.Valid {
		archivedAt = &sqlView.ArchivedAt.Time
	}

	return Weather{
		ID:                 sqlView.ID,
		CreatedAt:          sqlView.CreatedAt,
		UpdatedAt:          sqlView.UpdatedAt,
		ArchivedAt:         archivedAt,
		Latitude:           sqlView.Latitude,
		Longitude:          sqlView.Longitude,
		Temperature:        sqlView.Temperature,
		Humidity:           sqlView.Humidity,
		WindSpeed:          sqlView.WindSpeed,
		Clouds:             sqlView.Clouds,
		WeatherMain:        sqlView.WeatherMain,
		WeatherDesc:        sqlView.WeatherDesc,
		AgriculturalUnitId: sqlView.AgriculturalUnitId,
	}, nil
}

type WeatherStorage interface {
	InsertOrUpdate(weather Weather) error
}

type weatherStorage struct {
	querier storage.DBQuerier
	builder sq.StatementBuilderType
}

func NewWeatherStorage(querier storage.DBQuerier) WeatherStorage {
	return &weatherStorage{
		querier: querier,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *weatherStorage) InsertOrUpdate(w Weather) error {
	sqlView := WeatherToSqlView(w)

	builder := s.builder.Insert("weather").
		Columns(
			"id",
			"created_at",
			"updated_at",
			"archived_at",
			"latitude",
			"longitude",
			"temperature",
			"humidity",
			"wind_speed",
			"clouds",
			"weather_main",
			"weather_desc",
			"agricultural_unit_id",
		).
		Values(
			sqlView.ID,
			sqlView.CreatedAt,
			sqlView.UpdatedAt,
			sqlView.ArchivedAt,
			sqlView.Latitude,
			sqlView.Longitude,
			sqlView.Temperature,
			sqlView.Humidity,
			sqlView.WindSpeed,
			sqlView.Clouds,
			sqlView.WeatherMain,
			sqlView.WeatherDesc,
			sqlView.AgriculturalUnitId,
		).
		Suffix(`
            ON CONFLICT (id) DO UPDATE SET
                updated_at = EXCLUDED.updated_at,
                archived_at = EXCLUDED.archived_at,
                latitude = EXCLUDED.latitude,
                longitude = EXCLUDED.longitude,
                temperature = EXCLUDED.temperature,
                humidity = EXCLUDED.humidity,
                wind_speed = EXCLUDED.wind_speed,
                clouds = EXCLUDED.clouds,
                weather_main = EXCLUDED.weather_main,
                weather_desc = EXCLUDED.weather_desc,
                agricultural_unit_id = EXCLUDED.agricultural_unit_id
        `)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build InsertOrUpdate SQL for Weather: %w", err)
	}

	_, err = s.querier.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute InsertOrUpdate for Weather: %w", err)
	}

	return nil
}
