package agri_units

import (
	"agreste-ingestor/storage"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

type AgriculturalUnitSqlView struct {
	ID         string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ArchivedAt sql.NullTime

	IDNum     int
	Latitude  float64
	Longitude float64
}

func AgriculturalUnitToSqlView(unit AgriculturalUnit) AgriculturalUnitSqlView {
	sqlView := AgriculturalUnitSqlView{
		ID:        unit.ID.String(),
		CreatedAt: unit.CreatedAt,
		UpdatedAt: unit.UpdatedAt,
		IDNum:     unit.IDNum,
		Latitude:  unit.Latitude,
		Longitude: unit.Longitude,
	}

	if unit.ArchivedAt != nil {
		sqlView.ArchivedAt = sql.NullTime{Time: *unit.ArchivedAt, Valid: true}
	} else {
		sqlView.ArchivedAt = sql.NullTime{Valid: false}
	}

	return sqlView
}

func AgriculturalUnitFromSqlView(sqlView AgriculturalUnitSqlView) (AgriculturalUnit, error) {
	parsedID, err := uuid.Parse(sqlView.ID)
	if err != nil {
		return AgriculturalUnit{}, fmt.Errorf("failed to parse UUID from SQL view '%s': %w", sqlView.ID, err)
	}

	unit := AgriculturalUnit{
		ID:        parsedID,
		CreatedAt: sqlView.CreatedAt,
		UpdatedAt: sqlView.UpdatedAt,
		IDNum:     sqlView.IDNum,
		Latitude:  sqlView.Latitude,
		Longitude: sqlView.Longitude,
	}

	if sqlView.ArchivedAt.Valid {
		unit.ArchivedAt = &sqlView.ArchivedAt.Time
	} else {
		unit.ArchivedAt = nil
	}

	return unit, nil
}

type AgriUnitStorage interface {
	SelectAll() ([]AgriculturalUnit, error)
	InsertOrUpdate(unit AgriculturalUnit) error
}

type agriUnitStorage struct {
	querier storage.DBQuerier
	builder sq.StatementBuilderType
}

func NewAgriUnitStorage(querier storage.DBQuerier) AgriUnitStorage {

	return &agriUnitStorage{
		querier: querier,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *agriUnitStorage) SelectAll() ([]AgriculturalUnit, error) {
	queryBuilder := s.builder.Select(
		"id",
		"created_at",
		"updated_at",
		"archived_at",
		"id_num",
		"latitude",
		"longitude",
	).From("agricultural_units")

	sqlQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query with squirrel: %w", err)
	}

	rows, err := s.querier.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SelectAll query: %w", err)
	}
	defer rows.Close()

	var units []AgriculturalUnit
	for rows.Next() {
		var sqlView AgriculturalUnitSqlView
		err := rows.Scan(
			&sqlView.ID,
			&sqlView.CreatedAt,
			&sqlView.UpdatedAt,
			&sqlView.ArchivedAt,
			&sqlView.IDNum,
			&sqlView.Latitude,
			&sqlView.Longitude,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan agricultural unit row: %w", err)
		}

		domainUnit, err := AgriculturalUnitFromSqlView(sqlView)
		if err != nil {
			return nil, fmt.Errorf("failed to convert SQL view to domain model: %w", err)
		}
		units = append(units, domainUnit)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return units, nil
}

func (s *agriUnitStorage) InsertOrUpdate(unit AgriculturalUnit) error {
	sqlView := AgriculturalUnitToSqlView(unit)

	builder := s.builder.Insert("agricultural_units").
		Columns(
			"id",
			"created_at",
			"updated_at",
			"archived_at",
			"id_num",
			"latitude",
			"longitude",
		).
		Values(
			sqlView.ID,
			sqlView.CreatedAt,
			sqlView.UpdatedAt,
			sqlView.ArchivedAt,
			sqlView.IDNum,
			sqlView.Latitude,
			sqlView.Longitude,
		).
		Suffix(`
			ON CONFLICT (id) DO UPDATE SET
				updated_at = EXCLUDED.updated_at,
				archived_at = EXCLUDED.archived_at,
				id_num = EXCLUDED.id_num,
				latitude = EXCLUDED.latitude,
				longitude = EXCLUDED.longitude
		`)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build InsertOrUpdate SQL: %w", err)
	}

	_, err = s.querier.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute InsertOrUpdate: %w", err)
	}

	return nil
}
