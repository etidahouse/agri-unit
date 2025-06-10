package agri_units

import (
	"time"

	"github.com/google/uuid"
)

type AgriculturalUnitSurvey struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ArchivedAt *time.Time
	IDNum      int
	Year       int
	Data       map[string]interface{}
}

type AgriculturalUnitSurveyValue struct {
	IDNum int                    `json:"idNum"`
	Year  int                    `json:"Year"`
	Data  map[string]interface{} `json:"data"`
}

func CreateAgriculturalUnitSurvey(value AgriculturalUnitSurveyValue) AgriculturalUnitSurvey {
	now := time.Now()
	return AgriculturalUnitSurvey{
		ID:         uuid.New(),
		CreatedAt:  now,
		UpdatedAt:  now,
		IDNum:      value.IDNum,
		Year:       value.Year,
		Data:       value.Data,
		ArchivedAt: nil,
	}
}
