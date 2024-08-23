package userRepository

import (
	userModels "fiber-crud/internal/domain/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAll() ([]userModels.User, error)
	GetByID(id uuid.UUID) (userModels.User, error)
	GetByUsername(username string) (*userModels.User, error)
	GetByEmail(email string) (*userModels.User, error)
	Create(user userModels.User) (*userModels.User, error)
	Update(user userModels.User) error
	FindGoogleId(googleId string) (*userModels.User, error)
	Delete(id uuid.UUID) error
	Search(query string) ([]userModels.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAll() ([]userModels.User, error) {
	var users []userModels.User
	if err := r.db.Select("id, name, email, created_at").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetByID(id uuid.UUID) (userModels.User, error) {
	var u userModels.User
	if err := r.db.First(&u, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return userModels.User{}, nil
		}
		return userModels.User{}, err
	}
	return u, nil
}

func (r *userRepository) Search(query string) ([]userModels.User, error) {
	var users []userModels.User
	if err := r.db.Where("name LIKE ?", "%"+query+"%").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetByUsername(username string) (*userModels.User, error) {
	var u userModels.User
	if err := r.db.Where("name = ?", username).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByEmail(email string) (*userModels.User, error) {
	var u userModels.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Create(u userModels.User) (*userModels.User, error) {
	if err := r.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Update(u userModels.User) error {
	if err := r.db.Save(&u).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&userModels.User{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindGoogleId(googleId string) (*userModels.User, error) {
	var u userModels.User
	if err := r.db.Where("google_id = ?", googleId).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
