package agri_units

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateAgriculturalUnitSurvey(t *testing.T) {
	testValue := AgriculturalUnitSurveyValue{
		IDNum: 12345,
		Year:  2024,
		Data: map[string]interface{}{
			"farm_size_ha":  150.5,
			"crop_type":     "Wheat",
			"has_livestock": true,
			"location":      "France",
			"employees":     5,
		},
	}

	newSurvey := CreateAgriculturalUnitSurvey(testValue)

	now := time.Now()

	if newSurvey.ID == uuid.Nil {
		t.Errorf("Expected a non-nil UUID for ID, but got %v", newSurvey.ID)
	}

	if newSurvey.CreatedAt.IsZero() {
		t.Errorf("Expected CreatedAt to be set, but it's zero")
	}
	if newSurvey.UpdatedAt.IsZero() {
		t.Errorf("Expected UpdatedAt to be set, but it's zero")
	}

	tolerance := 100 * time.Millisecond
	if newSurvey.CreatedAt.Before(now.Add(-tolerance)) || newSurvey.CreatedAt.After(now.Add(tolerance)) {
		t.Errorf("CreatedAt mismatch: Expected around %v, got %v", now, newSurvey.CreatedAt)
	}
	if newSurvey.UpdatedAt.Before(now.Add(-tolerance)) || newSurvey.UpdatedAt.After(now.Add(tolerance)) {
		t.Errorf("UpdatedAt mismatch: Expected around %v, got %v", now, newSurvey.UpdatedAt)
	}
	if !newSurvey.CreatedAt.Equal(newSurvey.UpdatedAt) {
		t.Errorf("Expected CreatedAt and UpdatedAt to be equal for a new record, but got %v and %v", newSurvey.CreatedAt, newSurvey.UpdatedAt)
	}

	if newSurvey.IDNum != testValue.IDNum {
		t.Errorf("IDNum mismatch: Expected %d, got %d", testValue.IDNum, newSurvey.IDNum)
	}

	if newSurvey.Year != testValue.Year {
		t.Errorf("Year mismatch: Expected %d, got %d", testValue.Year, newSurvey.Year)
	}

	if newSurvey.Data == nil {
		t.Fatalf("Expected Data map to be initialized, but it's nil")
	}

	if !reflect.DeepEqual(newSurvey.Data, testValue.Data) {
		t.Errorf("Data map content mismatch: Expected %v, got %v", testValue.Data, newSurvey.Data)
	}

	if newSurvey.ArchivedAt != nil {
		t.Errorf("Expected ArchivedAt to be nil for a new record, but got %v", *newSurvey.ArchivedAt)
	}
}
