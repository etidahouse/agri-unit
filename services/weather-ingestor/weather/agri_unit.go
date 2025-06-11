package weather

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
