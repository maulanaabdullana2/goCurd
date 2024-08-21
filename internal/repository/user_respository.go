package userRepository

import (
	userModels "fiber-crud/internal/domain/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository mendefinisikan antarmuka untuk operasi CRUD pengguna
type UserRepository interface {
	GetAll() ([]userModels.User, error)
	GetByID(id uuid.UUID) (userModels.User, error)
	GetByUsername(username string) (*userModels.User, error)
	GetByEmail(email string) (*userModels.User, error)
	Create(user userModels.User) (*userModels.User, error)
	Update(user userModels.User) error
	Delete(id uuid.UUID) error
	Search(query string) ([]userModels.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// GetAll mengambil semua pengguna
func (r *userRepository) GetAll() ([]userModels.User, error) {
	var users []userModels.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetByID mengambil pengguna berdasarkan ID
func (r *userRepository) GetByID(id uuid.UUID) (userModels.User, error) {
	var u userModels.User
	if err := r.db.First(&u, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return userModels.User{}, nil // Return empty user if not found
		}
		return userModels.User{}, err
	}
	return u, nil
}

// Search mencari pengguna berdasarkan query string
func (r *userRepository) Search(query string) ([]userModels.User, error) {
	var users []userModels.User
	if err := r.db.Where("name LIKE ?", "%"+query+"%").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetByUsername mengambil pengguna berdasarkan username
func (r *userRepository) GetByUsername(username string) (*userModels.User, error) {
	var u userModels.User
	if err := r.db.Where("name = ?", username).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return &u, nil
}

// GetByEmail mengambil pengguna berdasarkan email
func (r *userRepository) GetByEmail(email string) (*userModels.User, error) {
	var u userModels.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return &u, nil
}

// Create membuat pengguna baru
func (r *userRepository) Create(u userModels.User) (*userModels.User, error) {
	if err := r.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// Update memperbarui pengguna yang ada
func (r *userRepository) Update(u userModels.User) error {
	if err := r.db.Save(&u).Error; err != nil {
		return err
	}
	return nil
}

// Delete menghapus pengguna berdasarkan ID
func (r *userRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&userModels.User{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
