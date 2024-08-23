package ProductModels

import (
	CommentModels "fiber-crud/internal/domain/comment"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Name        string
	Description string
	Price       float64
	ImageURL    string
	Comments    []CommentModels.Comment `gorm:"foreignKey:ProductID"`
	CreatedAt   time.Time
}
