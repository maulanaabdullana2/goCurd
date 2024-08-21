package ProductModels

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
}
