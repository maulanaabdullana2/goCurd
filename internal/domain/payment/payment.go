package paymentModels

import (
	"time"

	"github.com/google/uuid"
)

type PaymentModels struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	OrderID   string    `gorm:"type:uuid;not null"`
	Amount    int
	Status    string
	CreatedAt time.Time
}
