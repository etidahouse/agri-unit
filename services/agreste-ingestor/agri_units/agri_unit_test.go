package agri_units

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateAgriculturalUnit(t *testing.T) {

	testValue := AgriculturalUnitValue{
		IDNum:     123,
		Latitude:  48.8566,
		Longitude: 2.3522,
	}

	unit := CreateAgriculturalUnit(testValue)

	if unit.ID == uuid.Nil {
		t.Errorf("ID should not be nil")
	}

	if unit.CreatedAt.IsZero() {
		t.Errorf("CreatedAt should not be zero")
	}
	if unit.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt should not be zero")
	}

	if time.Since(unit.CreatedAt) > 5*time.Second {
		t.Errorf("CreatedAt is too old. Expected close to now, got %v", unit.CreatedAt)
	}
	if time.Since(unit.UpdatedAt) > 5*time.Second {
		t.Errorf("UpdatedAt is too old. Expected close to now, got %v", unit.UpdatedAt)
	}

	if !unit.CreatedAt.Equal(unit.UpdatedAt) {
		t.Errorf("CreatedAt and UpdatedAt should be equal upon creation. Got CreatedAt: %v, UpdatedAt: %v", unit.CreatedAt, unit.UpdatedAt)
	}

	if unit.IDNum != testValue.IDNum {
		t.Errorf("IDNum mismatch. Expected %d, got %d", testValue.IDNum, unit.IDNum)
	}
	if unit.Latitude != testValue.Latitude {
		t.Errorf("Latitude mismatch. Expected %f, got %f", testValue.Latitude, unit.Latitude)
	}
	if unit.Longitude != testValue.Longitude {
		t.Errorf("Longitude mismatch. Expected %f, got %f", testValue.Longitude, unit.Longitude)
	}

	if unit.ArchivedAt != nil {
		t.Errorf("ArchivedAt should be nil upon creation, got %v", unit.ArchivedAt)
	}
}
