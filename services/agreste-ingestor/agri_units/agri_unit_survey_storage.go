package agri_units

import (
	"agreste-ingestor/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

type AgriculturalUnitSurveySqlView struct {
	ID         string       `db:"id"`
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  time.Time    `db:"updated_at"`
	ArchivedAt sql.NullTime `db:"archived_at"`
	IDNum      int          `db:"id_num"`
	Year       int          `db:"year"`
	Data       []byte       `db:"data"`
}

func AgriculturalUnitSurveyToSqlView(survey AgriculturalUnitSurvey) (AgriculturalUnitSurveySqlView, error) {
	sqlView := AgriculturalUnitSurveySqlView{
		ID:        survey.ID.String(),
		CreatedAt: survey.CreatedAt,
		UpdatedAt: survey.UpdatedAt,
		IDNum:     survey.IDNum,
		Year:      survey.Year,
	}

	if survey.ArchivedAt != nil {
		sqlView.ArchivedAt = sql.NullTime{Time: *survey.ArchivedAt, Valid: true}
	} else {
		sqlView.ArchivedAt = sql.NullTime{Valid: false}
	}

	jsonData, err := json.Marshal(survey.Data)
	if err != nil {
		return AgriculturalUnitSurveySqlView{}, fmt.Errorf("failed to marshal survey data to JSON for SQL view: %w", err)
	}
	sqlView.Data = jsonData

	return sqlView, nil
}

func AgriculturalUnitSurveyFromSqlView(sqlView AgriculturalUnitSurveySqlView) (AgriculturalUnitSurvey, error) {
	parsedID, err := uuid.Parse(sqlView.ID)
	if err != nil {
		return AgriculturalUnitSurvey{}, fmt.Errorf("failed to parse UUID from SQL view '%s': %w", sqlView.ID, err)
	}

	survey := AgriculturalUnitSurvey{
		ID:        parsedID,
		CreatedAt: sqlView.CreatedAt,
		UpdatedAt: sqlView.UpdatedAt,
		IDNum:     sqlView.IDNum,
		Year:      sqlView.Year,
	}

	if sqlView.ArchivedAt.Valid {
		survey.ArchivedAt = &sqlView.ArchivedAt.Time
	} else {
		survey.ArchivedAt = nil
	}

	var dataMap map[string]interface{}
	if len(sqlView.Data) > 0 {
		err = json.Unmarshal(sqlView.Data, &dataMap)
		if err != nil {
			return AgriculturalUnitSurvey{}, fmt.Errorf("failed to unmarshal JSON data from SQL view: %w", err)
		}
	}
	survey.Data = dataMap

	return survey, nil
}

type AgriculturalUnitSurveyStorage interface {
	InsertOrUpdate(survey AgriculturalUnitSurvey) error
	SelectAll() ([]AgriculturalUnitSurvey, error)
}

type agriculturalUnitSurveyStorage struct {
	querier storage.DBQuerier
	builder sq.StatementBuilderType
}

func NewAgriculturalUnitSurveyStorage(querier storage.DBQuerier) AgriculturalUnitSurveyStorage {
	return &agriculturalUnitSurveyStorage{
		querier: querier,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *agriculturalUnitSurveyStorage) InsertOrUpdate(survey AgriculturalUnitSurvey) error {
	sqlView, err := AgriculturalUnitSurveyToSqlView(survey)
	if err != nil {
		return fmt.Errorf("failed to convert domain model to SQL view: %w", err)
	}

	builder := s.builder.Insert("agricultural_unit_surveys").
		Columns(
			"id",
			"created_at",
			"updated_at",
			"archived_at",
			"id_num",
			"year",
			"data",
		).
		Values(
			sqlView.ID,
			sqlView.CreatedAt,
			sqlView.UpdatedAt,
			sqlView.ArchivedAt,
			sqlView.IDNum,
			sqlView.Year,
			sqlView.Data,
		).
		Suffix(`
            ON CONFLICT (id) DO UPDATE SET
                updated_at = EXCLUDED.updated_at,
                archived_at = EXCLUDED.archived_at,
                id_num = EXCLUDED.id_num,
                year = EXCLUDED.year,
                data = EXCLUDED.data
        `)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build InsertOrUpdate SQL for AgriculturalUnitSurvey: %w", err)
	}

	_, err = s.querier.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute InsertOrUpdate for AgriculturalUnitSurvey: %w", err)
	}

	return nil
}

func (s *agriculturalUnitSurveyStorage) SelectAll() ([]AgriculturalUnitSurvey, error) {
	queryBuilder := s.builder.Select(
		"id",
		"created_at",
		"updated_at",
		"archived_at",
		"id_num",
		"year",
		"data",
	).From("agricultural_unit_surveys")

	sqlQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query with squirrel: %w", err)
	}

	rows, err := s.querier.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SelectAll query: %w", err)
	}
	defer rows.Close()

	var surveys []AgriculturalUnitSurvey
	for rows.Next() {
		var sqlView AgriculturalUnitSurveySqlView
		err := rows.Scan(
			&sqlView.ID,
			&sqlView.CreatedAt,
			&sqlView.UpdatedAt,
			&sqlView.ArchivedAt,
			&sqlView.IDNum,
			&sqlView.Year,
			&sqlView.Data,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan agricultural unit survey row: %w", err)
		}

		domainSurvey, err := AgriculturalUnitSurveyFromSqlView(sqlView)
		if err != nil {
			return nil, fmt.Errorf("failed to convert SQL view to domain model: %w", err)
		}
		surveys = append(surveys, domainSurvey)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return surveys, nil
}
