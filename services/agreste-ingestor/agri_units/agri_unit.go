package agri_units

import (
	"time"

	"github.com/google/uuid"
)

type AgriculturalUnit struct {
	ID         uuid.UUID  `json:"id"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	ArchivedAt *time.Time `json:"archivedAt,omitempty"`

	IDNum     int     `json:"idNum"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type AgriculturalUnitValue struct {
	IDNum     int     `json:"idNum"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func CreateAgriculturalUnit(value AgriculturalUnitValue) AgriculturalUnit {
	now := time.Now()
	return AgriculturalUnit{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		IDNum:     value.IDNum,
		Latitude:  value.Latitude,
		Longitude: value.Longitude,
	}
}
