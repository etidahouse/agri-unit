package agri_units

import (
	"agreste-ingestor/misc"
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestAgriculturalUnitSurveyToSqlView(t *testing.T) {
	t.Run("WithoutArchivedAt", func(t *testing.T) {
		now := time.Now().Truncate(time.Millisecond)
		domainSurvey := AgriculturalUnitSurvey{
			ID:         uuid.New(),
			CreatedAt:  now,
			UpdatedAt:  now,
			IDNum:      101,
			Year:       2023,
			Data:       map[string]interface{}{"crop": "corn", "area": 100.5},
			ArchivedAt: nil,
		}

		sqlView, err := AgriculturalUnitSurveyToSqlView(domainSurvey)
		if err != nil {
			t.Fatalf("AgriculturalUnitSurveyToSqlView returned an unexpected error: %v", err)
		}

		if sqlView.ID != domainSurvey.ID.String() {
			t.Errorf("ID mismatch. Expected %s, got %s", domainSurvey.ID.String(), sqlView.ID)
		}
		if !sqlView.CreatedAt.Equal(domainSurvey.CreatedAt) {
			t.Errorf("CreatedAt mismatch. Expected %v, got %v", domainSurvey.CreatedAt, sqlView.CreatedAt)
		}
		if !sqlView.UpdatedAt.Equal(domainSurvey.UpdatedAt) {
			t.Errorf("UpdatedAt mismatch. Expected %v, got %v", domainSurvey.UpdatedAt, sqlView.UpdatedAt)
		}
		if sqlView.ArchivedAt.Valid {
			t.Errorf("ArchivedAt.Valid should be false, got true")
		}
		if sqlView.IDNum != domainSurvey.IDNum {
			t.Errorf("IDNum mismatch. Expected %d, got %d", domainSurvey.IDNum, sqlView.IDNum)
		}
		if sqlView.Year != domainSurvey.Year {
			t.Errorf("Year mismatch. Expected %d, got %d", domainSurvey.Year, sqlView.Year)
		}

		expectedData, _ := json.Marshal(domainSurvey.Data)
		if string(sqlView.Data) != string(expectedData) {
			t.Errorf("Data mismatch. Expected %s, got %s", string(expectedData), string(sqlView.Data))
		}
	})

	t.Run("WithArchivedAt", func(t *testing.T) {
		now := time.Now().Truncate(time.Millisecond)
		archiveTime := now.Add(-24 * time.Hour)
		domainSurvey := AgriculturalUnitSurvey{
			ID:         uuid.New(),
			CreatedAt:  now.Add(-48 * time.Hour),
			UpdatedAt:  now,
			IDNum:      202,
			Year:       2024,
			Data:       map[string]interface{}{"crop": "wheat", "yield": 500.0},
			ArchivedAt: &archiveTime,
		}

		sqlView, err := AgriculturalUnitSurveyToSqlView(domainSurvey)
		if err != nil {
			t.Fatalf("AgriculturalUnitSurveyToSqlView returned an unexpected error: %v", err)
		}

		if !sqlView.ArchivedAt.Valid {
			t.Errorf("ArchivedAt.Valid should be true, got false")
		}
		if !sqlView.ArchivedAt.Time.Equal(*domainSurvey.ArchivedAt) {
			t.Errorf("ArchivedAt.Time mismatch. Expected %v, got %v", *domainSurvey.ArchivedAt, sqlView.ArchivedAt.Time)
		}
	})

	t.Run("MarshalError", func(t *testing.T) {
		domainSurvey := AgriculturalUnitSurvey{
			ID:         uuid.New(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			IDNum:      1,
			Year:       2023,
			Data:       map[string]interface{}{"unmarshalable": make(chan int)},
			ArchivedAt: nil,
		}

		_, err := AgriculturalUnitSurveyToSqlView(domainSurvey)
		if err == nil {
			t.Error("AgriculturalUnitSurveyToSqlView expected an error for unmarshalable data, but got none.")
		}

		var jsonTypeError *json.UnsupportedTypeError
		if !errors.As(err, &jsonTypeError) && err.Error() != "failed to marshal survey data to JSON for SQL view: json: unsupported type: chan int" {
			t.Errorf("AgriculturalUnitSurveyToSqlView expected a JSON marshal error, got %v", err)
		}
	})
}

func TestAgriculturalUnitSurveyFromSqlView(t *testing.T) {
	t.Run("FromSqlViewWithoutArchivedAt", func(t *testing.T) {
		now := time.Now().Truncate(time.Millisecond)
		id := uuid.New()
		dataJSON, _ := json.Marshal(map[string]interface{}{"field_a": "value_a", "field_b": 123})
		sqlView := AgriculturalUnitSurveySqlView{
			ID:         id.String(),
			CreatedAt:  now.Add(-72 * time.Hour),
			UpdatedAt:  now.Add(-24 * time.Hour),
			ArchivedAt: sql.NullTime{Valid: false},
			IDNum:      303,
			Year:       2025,
			Data:       dataJSON,
		}

		domainSurvey, err := AgriculturalUnitSurveyFromSqlView(sqlView)
		if err != nil {
			t.Fatalf("AgriculturalUnitSurveyFromSqlView returned an unexpected error: %v", err)
		}

		if domainSurvey.ID != id {
			t.Errorf("ID mismatch. Expected %s, got %s", id.String(), domainSurvey.ID.String())
		}
		if !domainSurvey.CreatedAt.Equal(sqlView.CreatedAt) {
			t.Errorf("CreatedAt mismatch. Expected %v, got %v", sqlView.CreatedAt, domainSurvey.CreatedAt)
		}
		if !domainSurvey.UpdatedAt.Equal(sqlView.UpdatedAt) {
			t.Errorf("UpdatedAt mismatch. Expected %v, got %v", sqlView.UpdatedAt, domainSurvey.UpdatedAt)
		}
		if domainSurvey.ArchivedAt != nil {
			t.Errorf("ArchivedAt should be nil, got %v", domainSurvey.ArchivedAt)
		}
		if domainSurvey.IDNum != sqlView.IDNum {
			t.Errorf("IDNum mismatch. Expected %d, got %d", sqlView.IDNum, domainSurvey.IDNum)
		}
		if domainSurvey.Year != sqlView.Year {
			t.Errorf("Year mismatch. Expected %d, got %d", sqlView.Year, domainSurvey.Year)
		}
		expectedData := map[string]interface{}{"field_a": "value_a", "field_b": float64(123)}
		if !reflect.DeepEqual(domainSurvey.Data, expectedData) {
			t.Errorf("Data mismatch. Expected %v, got %v", expectedData, domainSurvey.Data)
		}
	})

	t.Run("FromSqlViewWithArchivedAt", func(t *testing.T) {
		now := time.Now().Truncate(time.Millisecond)
		id := uuid.New()
		archiveTime := now.Add(-12 * time.Hour)
		dataJSON, _ := json.Marshal(map[string]interface{}{"status": "archived"})
		sqlView := AgriculturalUnitSurveySqlView{
			ID:         id.String(),
			CreatedAt:  now.Add(-36 * time.Hour),
			UpdatedAt:  now,
			ArchivedAt: sql.NullTime{Time: archiveTime, Valid: true},
			IDNum:      404,
			Year:       2026,
			Data:       dataJSON,
		}

		domainSurvey, err := AgriculturalUnitSurveyFromSqlView(sqlView)
		if err != nil {
			t.Fatalf("AgriculturalUnitSurveyFromSqlView returned an unexpected error: %v", err)
		}

		if domainSurvey.ArchivedAt == nil {
			t.Errorf("ArchivedAt should not be nil")
		}
		if !domainSurvey.ArchivedAt.Equal(sqlView.ArchivedAt.Time) {
			t.Errorf("ArchivedAt mismatch. Expected %v, got %v", sqlView.ArchivedAt.Time, *domainSurvey.ArchivedAt)
		}
	})

	t.Run("FromSqlViewWithMalformedJSON", func(t *testing.T) {
		now := time.Now().Truncate(time.Millisecond)
		id := uuid.New()
		sqlView := AgriculturalUnitSurveySqlView{
			ID:         id.String(),
			CreatedAt:  now,
			UpdatedAt:  now,
			IDNum:      606,
			Year:       2028,
			Data:       []byte(`{"key": "value", "another_key": `),
			ArchivedAt: sql.NullTime{Valid: false},
		}

		_, err := AgriculturalUnitSurveyFromSqlView(sqlView)
		if err == nil {
			t.Error("AgriculturalUnitSurveyFromSqlView expected an error for malformed JSON, but got none.")
		}
		var syntaxError *json.SyntaxError
		if !errors.As(err, &syntaxError) {
			t.Errorf("AgriculturalUnitSurveyFromSqlView expected a JSON syntax error, got %v", err)
		}
	})
}

func TestAgriculturalUnitSurveyStorage_SelectAll(t *testing.T) {
	t.Run("SuccessfulSelectAll", func(t *testing.T) {
		mockQuerierInstance, sqlMock, err := misc.NewMockQuerier(t)
		if err != nil {
			t.Fatalf("failed to create mock querier: %v", err)
		}
		defer mockQuerierInstance.Db.Close()

		storage := NewAgriculturalUnitSurveyStorage(mockQuerierInstance)

		id1 := uuid.New()
		createdAt1 := time.Now().Add(-24 * time.Hour).Truncate(time.Millisecond)
		updatedAt1 := time.Now().Truncate(time.Millisecond)
		archivedAt1 := sql.NullTime{Time: time.Now().Add(-12 * time.Hour).Truncate(time.Millisecond), Valid: true}
		data1 := map[string]interface{}{"crop": "corn", "area": 100.5}
		dataJSON1, _ := json.Marshal(data1)

		id2 := uuid.New()
		createdAt2 := time.Now().Add(-48 * time.Hour).Truncate(time.Millisecond)
		updatedAt2 := time.Now().Add(-2 * time.Hour).Truncate(time.Millisecond)
		archivedAt2 := sql.NullTime{Valid: false}
		data2 := map[string]interface{}{"region": "south", "yield": 50.2}
		dataJSON2, _ := json.Marshal(data2)

		rows := sqlMock.NewRows([]string{"id", "created_at", "updated_at", "archived_at", "id_num", "year", "data"}).
			AddRow(id1.String(), createdAt1, updatedAt1, archivedAt1, 1, 2023, dataJSON1).
			AddRow(id2.String(), createdAt2, updatedAt2, archivedAt2, 2, 2024, dataJSON2)

		sqlMock.ExpectQuery("SELECT id, created_at, updated_at, archived_at, id_num, year, data FROM agricultural_unit_surveys").
			WillReturnRows(rows)

		surveys, err := storage.SelectAll()
		if err != nil {
			t.Fatalf("SelectAll returned an unexpected error: %v", err)
		}

		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

		if len(surveys) != 2 {
			t.Fatalf("Expected 2 surveys, got %d", len(surveys))
		}

		if surveys[0].ID != id1 {
			t.Errorf("Survey 1 ID mismatch. Expected %s, got %s", id1, surveys[0].ID)
		}
		if !surveys[0].CreatedAt.Equal(createdAt1) {
			t.Errorf("Survey 1 CreatedAt mismatch. Expected %v, got %v", createdAt1, surveys[0].CreatedAt)
		}
		if !surveys[0].UpdatedAt.Equal(updatedAt1) {
			t.Errorf("Survey 1 UpdatedAt mismatch. Expected %v, got %v", updatedAt1, surveys[0].UpdatedAt)
		}
		if surveys[0].ArchivedAt == nil || !surveys[0].ArchivedAt.Equal(archivedAt1.Time) {
			t.Errorf("Survey 1 ArchivedAt mismatch. Expected %v, got %v", archivedAt1.Time, surveys[0].ArchivedAt)
		}
		if surveys[0].IDNum != 1 {
			t.Errorf("Survey 1 IDNum mismatch. Expected %d, got %d", 1, surveys[0].IDNum)
		}
		if surveys[0].Year != 2023 {
			t.Errorf("Survey 1 Year mismatch. Expected %d, got %d", 2023, surveys[0].Year)
		}

		expectedData1 := map[string]interface{}{"crop": "corn", "area": float64(100.5)}
		if !reflect.DeepEqual(surveys[0].Data, expectedData1) {
			t.Errorf("Survey 1 Data mismatch. Expected %v, got %v", expectedData1, surveys[0].Data)
		}

		if surveys[1].ID != id2 {
			t.Errorf("Survey 2 ID mismatch. Expected %s, got %s", id2, surveys[1].ID)
		}
		if surveys[1].ArchivedAt != nil {
			t.Errorf("Survey 2 ArchivedAt should be nil, got %v", surveys[1].ArchivedAt)
		}
		if surveys[1].IDNum != 2 {
			t.Errorf("Survey 2 IDNum mismatch. Expected %d, got %d", 2, surveys[1].IDNum)
		}
		if surveys[1].Year != 2024 {
			t.Errorf("Survey 2 Year mismatch. Expected %d, got %d", 2024, surveys[1].Year)
		}
		expectedData2 := map[string]interface{}{"region": "south", "yield": float64(50.2)}
		if !reflect.DeepEqual(surveys[1].Data, expectedData2) {
			t.Errorf("Survey 2 Data mismatch. Expected %v, got %v", expectedData2, surveys[1].Data)
		}
	})

	t.Run("QueryError", func(t *testing.T) {
		mockQuerierInstance, sqlMock, err := misc.NewMockQuerier(t)
		if err != nil {
			t.Fatalf("Failed to create mock querier: %v", err)
		}
		defer mockQuerierInstance.Db.Close()

		expectedErr := errors.New("database query failed")
		sqlMock.ExpectQuery("SELECT id, created_at, updated_at, archived_at, id_num, year, data FROM agricultural_unit_surveys").
			WillReturnError(expectedErr)

		storage := NewAgriculturalUnitSurveyStorage(mockQuerierInstance)

		_, err = storage.SelectAll()
		if err == nil {
			t.Error("SelectAll expected an error, but got none")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}

		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})

	t.Run("RowsError", func(t *testing.T) {
		mockQuerierInstance, sqlMock, err := misc.NewMockQuerier(t)
		if err != nil {
			t.Fatalf("Failed to create mock querier: %v", err)
		}
		defer mockQuerierInstance.Db.Close()

		expectedErr := errors.New("rows iteration error")
		rows := sqlMock.NewRows([]string{"id", "created_at", "updated_at", "archived_at", "id_num", "year", "data"}).
			AddRow(uuid.New().String(), time.Now(), time.Now(), sql.NullTime{Valid: false}, 1, 2023, []byte("{}")).
			RowError(0, expectedErr)

		sqlMock.ExpectQuery("SELECT id, created_at, updated_at, archived_at, id_num, year, data FROM agricultural_unit_surveys").
			WillReturnRows(rows)

		storage := NewAgriculturalUnitSurveyStorage(mockQuerierInstance)

		_, err = storage.SelectAll()
		if err == nil {
			t.Error("SelectAll expected a rows.Err(), but got none")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}

		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})
}

func TestAgriculturalUnitSurveyStorage_InsertOrUpdate(t *testing.T) {

	t.Run("SuccessfulInsertOrUpdate", func(t *testing.T) {
		mockQuerierInstance, sqlMock, err := misc.NewMockQuerier(t)
		if err != nil {
			t.Fatalf("Failed to create mock querier: %v", err)
		}
		defer mockQuerierInstance.Db.Close()

		storage := NewAgriculturalUnitSurveyStorage(mockQuerierInstance)

		surveyID := uuid.New()
		now := time.Now().Truncate(time.Millisecond)
		archivedAt := now.Add(-time.Hour).Truncate(time.Millisecond)
		surveyData := map[string]interface{}{
			"type":      "vegetable",
			"cultivars": []string{"tomato", "cucumber"},
			"soil_ph":   float64(6.5),
		}

		domainSurvey := AgriculturalUnitSurvey{
			ID:         surveyID,
			CreatedAt:  now,
			UpdatedAt:  now,
			ArchivedAt: &archivedAt,
			IDNum:      123,
			Year:       2025,
			Data:       surveyData,
		}

		sqlView, err := AgriculturalUnitSurveyToSqlView(domainSurvey)
		if err != nil {
			t.Fatalf("Failed to convert domain survey to SQL view: %v", err)
		}

		sqlMock.ExpectExec(
			regexp.QuoteMeta(`INSERT INTO agricultural_unit_surveys (id,created_at,updated_at,archived_at,id_num,year,data) VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (id) DO UPDATE SET updated_at = EXCLUDED.updated_at, archived_at = EXCLUDED.archived_at, id_num = EXCLUDED.id_num, year = EXCLUDED.year, data = EXCLUDED.data`),
		).WithArgs(
			sqlView.ID,
			sqlView.CreatedAt,
			sqlView.UpdatedAt,
			sqlView.ArchivedAt,
			sqlView.IDNum,
			sqlView.Year,
			sqlView.Data,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err = storage.InsertOrUpdate(domainSurvey)
		if err != nil {
			t.Fatalf("InsertOrUpdate returned an unexpected error: %v", err)
		}

		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})
}
