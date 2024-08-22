package userModels

import (
	"time"

	"github.com/google/uuid"
)

// User mewakili tabel pengguna dalam database
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	Avatar    string    `json:"avatar"`
	GoogleID  string    `json:"google_id"`
	CreatedAt time.Time `json:"created_at"`
}
