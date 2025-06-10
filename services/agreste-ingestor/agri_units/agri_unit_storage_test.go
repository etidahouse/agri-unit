package agri_units

import (
	"agreste-ingestor/misc"
	"database/sql"
	"database/sql/driver"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	go_sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestAgriculturalUnitToSqlView(t *testing.T) {

	t.Run("WithoutArchivedAt", func(t *testing.T) {
		now := time.Now().Truncate(time.Millisecond)
		domainUnit := AgriculturalUnit{
			ID:         uuid.New(),
			CreatedAt:  now,
			UpdatedAt:  now,
			IDNum:      101,
			Latitude:   45.0,
			Longitude:  5.0,
			ArchivedAt: nil,
		}

		sqlView := AgriculturalUnitToSqlView(domainUnit)

		if sqlView.ID != domainUnit.ID.String() {
			t.Errorf("ID mismatch. Expected %s, got %s", domainUnit.ID.String(), sqlView.ID)
		}
		if !sqlView.CreatedAt.Equal(domainUnit.CreatedAt) {
			t.Errorf("CreatedAt mismatch. Expected %v, got %v", domainUnit.CreatedAt, sqlView.CreatedAt)
		}
		if !sqlView.UpdatedAt.Equal(domainUnit.UpdatedAt) {
			t.Errorf("UpdatedAt mismatch. Expected %v, got %v", domainUnit.UpdatedAt, sqlView.UpdatedAt)
		}
		if sqlView.ArchivedAt.Valid {
			t.Errorf("ArchivedAt.Valid should be false, got true")
		}
		if sqlView.IDNum != domainUnit.IDNum {
			t.Errorf("IDNum mismatch. Expected %d, got %d", domainUnit.IDNum, sqlView.IDNum)
		}
		if sqlView.Latitude != domainUnit.Latitude {
			t.Errorf("Latitude mismatch. Expected %f, got %f", domainUnit.Latitude, sqlView.Latitude)
		}
		if sqlView.Longitude != domainUnit.Longitude {
			t.Errorf("Longitude mismatch. Expected %f, got %f", domainUnit.Longitude, sqlView.Longitude)
		}
	})

	t.Run("WithArchivedAt", func(t *testing.T) {
		now := time.Now().Truncate(time.Millisecond)
		archiveTime := now.Add(-24 * time.Hour)
		domainUnit := AgriculturalUnit{
			ID:         uuid.New(),
			CreatedAt:  now.Add(-48 * time.Hour),
			UpdatedAt:  now,
			IDNum:      202,
			Latitude:   48.0,
			Longitude:  8.0,
			ArchivedAt: &archiveTime,
		}

		sqlView := AgriculturalUnitToSqlView(domainUnit)

		if !sqlView.ArchivedAt.Valid {
			t.Errorf("ArchivedAt.Valid should be true, got false")
		}
		if !sqlView.ArchivedAt.Time.Equal(*domainUnit.ArchivedAt) {
			t.Errorf("ArchivedAt.Time mismatch. Expected %v, got %v", *domainUnit.ArchivedAt, sqlView.ArchivedAt.Time)
		}
	})
}

func TestAgriculturalUnitFromSqlView(t *testing.T) {
	t.Run("FromSqlViewWithoutArchivedAt", func(t *testing.T) {
		now := time.Now().Truncate(time.Millisecond)
		id := uuid.New()
		sqlView := AgriculturalUnitSqlView{
			ID:         id.String(),
			CreatedAt:  now.Add(-72 * time.Hour),
			UpdatedAt:  now.Add(-24 * time.Hour),
			ArchivedAt: sql.NullTime{Valid: false},
			IDNum:      303,
			Latitude:   42.5,
			Longitude:  6.5,
		}

		domainUnit, err := AgriculturalUnitFromSqlView(sqlView)
		if err != nil {
			t.Fatalf("AgriculturalUnitFromSqlView returned an unexpected error: %v", err)
		}

		if domainUnit.ID != id {
			t.Errorf("ID mismatch. Expected %s, got %s", id.String(), domainUnit.ID.String())
		}
		if !domainUnit.CreatedAt.Equal(sqlView.CreatedAt) {
			t.Errorf("CreatedAt mismatch. Expected %v, got %v", sqlView.CreatedAt, domainUnit.CreatedAt)
		}
		if !domainUnit.UpdatedAt.Equal(sqlView.UpdatedAt) {
			t.Errorf("UpdatedAt mismatch. Expected %v, got %v", sqlView.UpdatedAt, domainUnit.UpdatedAt)
		}
		if domainUnit.ArchivedAt != nil {
			t.Errorf("ArchivedAt should be nil, got %v", *domainUnit.ArchivedAt)
		}
		if domainUnit.IDNum != sqlView.IDNum {
			t.Errorf("IDNum mismatch. Expected %d, got %d", sqlView.IDNum, domainUnit.IDNum)
		}
		if domainUnit.Latitude != sqlView.Latitude {
			t.Errorf("Latitude mismatch. Expected %f, got %f", sqlView.Latitude, domainUnit.Latitude)
		}
		if domainUnit.Longitude != sqlView.Longitude {
			t.Errorf("Longitude mismatch. Expected %f, got %f", sqlView.Longitude, domainUnit.Longitude)
		}
	})

	t.Run("FromSqlViewWithArchivedAt", func(t *testing.T) {
		now := time.Now().Truncate(time.Millisecond)
		id := uuid.New()
		archiveTime := now.Add(-12 * time.Hour)
		sqlView := AgriculturalUnitSqlView{
			ID:         id.String(),
			CreatedAt:  now.Add(-36 * time.Hour),
			UpdatedAt:  now,
			ArchivedAt: sql.NullTime{Time: archiveTime, Valid: true},
			IDNum:      404,
			Latitude:   49.1,
			Longitude:  0.5,
		}

		domainUnit, err := AgriculturalUnitFromSqlView(sqlView)
		if err != nil {
			t.Fatalf("AgriculturalUnitFromSqlView returned an unexpected error: %v", err)
		}

		if domainUnit.ArchivedAt == nil {
			t.Errorf("ArchivedAt should not be nil")
		}
		if !domainUnit.ArchivedAt.Equal(sqlView.ArchivedAt.Time) {
			t.Errorf("ArchivedAt mismatch. Expected %v, got %v", sqlView.ArchivedAt.Time, *domainUnit.ArchivedAt)
		}
	})

	t.Run("FromSqlViewWithInvalidUUID", func(t *testing.T) {
		now := time.Now().Truncate(time.Millisecond)
		sqlView := AgriculturalUnitSqlView{
			ID:         "invalid-uuid-string",
			CreatedAt:  now,
			UpdatedAt:  now,
			IDNum:      505,
			Latitude:   1.0,
			Longitude:  1.0,
			ArchivedAt: sql.NullTime{Valid: false},
		}

		_, err := AgriculturalUnitFromSqlView(sqlView)
		if err == nil {
			t.Errorf("AgriculturalUnitFromSqlView expected an error for invalid UUID, but got none.")
		}
	})
}

func TestSelectAll(t *testing.T) {

	mockQuerierInstance, sqlMock, err := misc.NewMockQuerier(t)
	if err != nil {
		t.Fatalf("failed to create mock querier: %v", err)
	}
	defer mockQuerierInstance.Db.Close()

	storage := NewAgriUnitStorage(mockQuerierInstance)

	id1 := uuid.New()
	id2 := uuid.New()
	now := time.Now().Truncate(time.Millisecond)

	unit1CreatedAt := now.Add(-time.Hour * 24 * 7)
	unit1UpdatedAt := now.Add(-time.Hour)

	unit1ArchivedAt := sql.NullTime{Valid: false}

	unit2CreatedAt := now.Add(-time.Hour * 24 * 30)
	unit2UpdatedAt := now.Add(-time.Hour * 5)
	unit2ArchiveTime := now.Add(-time.Hour * 2)

	unit2ArchivedAt := sql.NullTime{Time: unit2ArchiveTime, Valid: true}

	expectedUnits := []AgriculturalUnit{
		{
			ID:         id1,
			CreatedAt:  unit1CreatedAt,
			UpdatedAt:  unit1UpdatedAt,
			ArchivedAt: nil,
			IDNum:      1,
			Latitude:   45.0,
			Longitude:  5.0,
		},
		{
			ID:         id2,
			CreatedAt:  unit2CreatedAt,
			UpdatedAt:  unit2UpdatedAt,
			ArchivedAt: &unit2ArchiveTime,
			IDNum:      2,
			Latitude:   48.0,
			Longitude:  2.0,
		},
	}

	columns := []string{"id", "created_at", "updated_at", "archived_at", "id_num", "latitude", "longitude"}
	sqlMock.ExpectQuery(`SELECT id, created_at, updated_at, archived_at, id_num, latitude, longitude FROM agricultural_units`).
		WillReturnRows(sqlMock.NewRows(columns).
			AddRow(id1.String(), unit1CreatedAt, unit1UpdatedAt, unit1ArchivedAt, 1, 45.0, 5.0).
			AddRow(id2.String(), unit2CreatedAt, unit2UpdatedAt, unit2ArchivedAt, 2, 48.0, 2.0))

	retrievedUnits, err := storage.SelectAll()

	if err != nil {
		t.Fatalf("SelectAll() returned an unexpected error: %v", err)
	}

	if len(retrievedUnits) != len(expectedUnits) {
		t.Fatalf("expected %d units, got %d", len(expectedUnits), len(retrievedUnits))
	}

	for i := range retrievedUnits {
		if retrievedUnits[i].ID != expectedUnits[i].ID {
			t.Errorf("unit %d ID mismatch. Expected %s, got %s", i, expectedUnits[i].ID, retrievedUnits[i].ID)
		}
		if !retrievedUnits[i].CreatedAt.Equal(expectedUnits[i].CreatedAt) {
			t.Errorf("unit %d CreatedAt mismatch. Expected %v, got %v", i, expectedUnits[i].CreatedAt, retrievedUnits[i].CreatedAt)
		}
		if !retrievedUnits[i].UpdatedAt.Equal(expectedUnits[i].UpdatedAt) {
			t.Errorf("unit %d UpdatedAt mismatch. Expected %v, got %v", i, expectedUnits[i].UpdatedAt, retrievedUnits[i].UpdatedAt)
		}

		if (retrievedUnits[i].ArchivedAt == nil) != (expectedUnits[i].ArchivedAt == nil) {
			t.Errorf("unit %d ArchivedAt nil status mismatch. Retrieved nil: %t, Expected nil: %t",
				i, retrievedUnits[i].ArchivedAt == nil, expectedUnits[i].ArchivedAt == nil)
		}
		if retrievedUnits[i].ArchivedAt != nil && !retrievedUnits[i].ArchivedAt.Equal(*expectedUnits[i].ArchivedAt) {
			t.Errorf("unit %d ArchivedAt value mismatch. Expected %v, got %v", i, *expectedUnits[i].ArchivedAt, *retrievedUnits[i].ArchivedAt)
		}

		if retrievedUnits[i].IDNum != expectedUnits[i].IDNum {
			t.Errorf("unit %d IDNum mismatch. Expected %d, got %d", i, expectedUnits[i].IDNum, retrievedUnits[i].IDNum)
		}
		if retrievedUnits[i].Latitude != expectedUnits[i].Latitude {
			t.Errorf("unit %d Latitude mismatch. Expected %f, got %f", i, expectedUnits[i].Latitude, retrievedUnits[i].Latitude)
		}
		if retrievedUnits[i].Longitude != expectedUnits[i].Longitude {
			t.Errorf("unit %d Longitude mismatch. Expected %f, got %f", i, expectedUnits[i].Longitude, retrievedUnits[i].Longitude)
		}
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSelectAll_QueryError(t *testing.T) {
	mockQuerierInstance, sqlMock, err := misc.NewMockQuerier(t)
	if err != nil {
		t.Fatalf("failed to create mock querier: %v", err)
	}
	defer mockQuerierInstance.Db.Close()

	storage := NewAgriUnitStorage(mockQuerierInstance)
	expectedError := errors.New("simulated database query error")

	sqlMock.ExpectQuery(`SELECT id, created_at, updated_at, archived_at, id_num, latitude, longitude FROM agricultural_units`).
		WillReturnError(expectedError)

	_, err = storage.SelectAll()
	if err == nil {
		t.Error("SelectAll() expected an error, but got none")
	}
	if !errors.Is(err, expectedError) {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSelectAll_ScanError(t *testing.T) {
	mockQuerierInstance, sqlMock, err := misc.NewMockQuerier(t)
	if err != nil {
		t.Fatalf("failed to create mock querier: %v", err)
	}
	defer mockQuerierInstance.Db.Close()

	storage := NewAgriUnitStorage(mockQuerierInstance)

	columns := []string{"id", "created_at"}
	sqlMock.ExpectQuery(`SELECT id, created_at, updated_at, archived_at, id_num, latitude, longitude FROM agricultural_units`).
		WillReturnRows(sqlMock.NewRows(columns).
			AddRow(uuid.New().String(), time.Now()))

	_, err = storage.SelectAll()
	if err == nil {
		t.Error("SelectAll() expected a scan error, but got none")
	}

	if !strings.Contains(err.Error(), "scan") && !strings.Contains(err.Error(), "expected") {
		t.Errorf("expected scan error, got %v", err)
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertOrUpdate_Success(t *testing.T) {
	mockQuerierInstance, sqlMock, err := misc.NewMockQuerier(t)
	if err != nil {
		t.Fatalf("failed to create mock querier: %v", err)
	}
	defer mockQuerierInstance.Db.Close()

	storage := NewAgriUnitStorage(mockQuerierInstance)

	testID := uuid.New()
	testCreatedAt := time.Now().Add(-24 * time.Hour).Truncate(time.Millisecond)
	testUpdatedAt := time.Now().Truncate(time.Millisecond)
	testArchivedAt := time.Now().Add(-1 * time.Hour).Truncate(time.Millisecond)
	testUnit := AgriculturalUnit{
		ID:         testID,
		CreatedAt:  testCreatedAt,
		UpdatedAt:  testUpdatedAt,
		ArchivedAt: &testArchivedAt,
		IDNum:      12345,
		Latitude:   10.123,
		Longitude:  20.456,
	}

	expectedSQL := "INSERT INTO agricultural_units (id,created_at,updated_at,archived_at,id_num,latitude,longitude) VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (id) DO UPDATE SET updated_at = EXCLUDED.updated_at, archived_at = EXCLUDED.archived_at, id_num = EXCLUDED.id_num, latitude = EXCLUDED.latitude, longitude = EXCLUDED.longitude"

	expectedArgs := []interface{}{
		testUnit.ID.String(),
		testUnit.CreatedAt,
		testUnit.UpdatedAt,
		sql.NullTime{Time: *testUnit.ArchivedAt, Valid: true},
		testUnit.IDNum,
		testUnit.Latitude,
		testUnit.Longitude,
	}

	var driverArgs []driver.Value
	for _, arg := range expectedArgs {
		driverArgs = append(driverArgs, arg)
	}

	sqlMock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
		WithArgs(driverArgs...).
		WillReturnResult(go_sqlmock.NewResult(1, 1))

	err = storage.InsertOrUpdate(testUnit)

	if err != nil {
		t.Fatalf("InsertOrUpdate() returned an unexpected error: %v", err)
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
