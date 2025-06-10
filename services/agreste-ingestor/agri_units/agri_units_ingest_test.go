package agri_units

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

type MockAgriUnitStorage struct {
	Units []AgriculturalUnit
	Error error
}

func (m *MockAgriUnitStorage) SelectAll() ([]AgriculturalUnit, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Units, nil
}

func (m *MockAgriUnitStorage) InsertOrUpdate(unit AgriculturalUnit) error {
	if m.Error != nil {
		return m.Error
	}

	for i, u := range m.Units {
		if u.IDNum == unit.IDNum {
			m.Units[i] = unit
			return nil
		}
	}
	m.Units = append(m.Units, unit)
	return nil
}

type MockAgriculturalUnitSurveyStorage struct {
	Surveys []AgriculturalUnitSurvey
	Error   error
}

func (m *MockAgriculturalUnitSurveyStorage) SelectAll() ([]AgriculturalUnitSurvey, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Surveys, nil
}

func (m *MockAgriculturalUnitSurveyStorage) InsertOrUpdate(survey AgriculturalUnitSurvey) error {
	if m.Error != nil {
		return m.Error
	}

	for i, s := range m.Surveys {
		if s.IDNum == survey.IDNum && s.Year == survey.Year {
			m.Surveys[i] = survey
			return nil
		}
	}
	m.Surveys = append(m.Surveys, survey)
	return nil
}

var mockCSVData [][]string
var mockCSVError error

func MockDownloadZipAndReadSpecificCSV(endpointURL string, fileName string) ([][]string, error) {
	return mockCSVData, mockCSVError
}

func TestHandleAgriUnitSurveyIngest_Success(t *testing.T) {
	existingUnit1 := AgriculturalUnit{ID: uuid.New(), IDNum: 101, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	existingSurvey1 := AgriculturalUnitSurvey{ID: uuid.New(), IDNum: 101, Year: 2023, Data: map[string]interface{}{"foo": "bar"}, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	mockAgriUnitStorage := &MockAgriUnitStorage{
		Units: []AgriculturalUnit{existingUnit1},
	}
	mockAgriUnitSurveyStorage := &MockAgriculturalUnitSurveyStorage{
		Surveys: []AgriculturalUnitSurvey{existingSurvey1},
	}

	agresteZipURL := "https://agreste.agriculture.gouv.fr/agreste-web/download/service/SV-Accès micro données RICA/RicaMicrodonnées2023_v2.zip"
	err := HandleAgriUnitSurveyIngest(agresteZipURL, "Rica_France_micro_Donnees_ex2023.csv", mockAgriUnitStorage, mockAgriUnitSurveyStorage)

	if err != nil {
		t.Fatalf("HandleAgriUnitSurveyIngest failed with error: %v", err)
	}

}
