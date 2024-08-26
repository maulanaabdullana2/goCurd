package cartModels

import (
	ProductModels "fiber-crud/internal/domain/product"

	"github.com/google/uuid"
)

type CartModels struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Quantity  int
	Product   ProductModels.Product `gorm:"foreignKey:ProductID"`
}
