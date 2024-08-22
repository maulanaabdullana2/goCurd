package Userusecase

import (
	"errors"
	"time"

	userModels "fiber-crud/internal/domain/user"
	userRepository "fiber-crud/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound           = errors.New("entity not found")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrEmailTaken         = errors.New("email already taken")
	ErrUsernameValidate   = errors.New("username cannot be empty")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserUsecase interface {
	GetUsers() ([]userModels.User, error)
	GetUserByID(id uuid.UUID) (userModels.User, error)
	CreateUser(user userModels.User) (*userModels.User, error)
	UpdateUser(user userModels.User) error
	DeleteUser(id uuid.UUID) error
	SearchUsers(query string) ([]userModels.User, error)
	LoginOrSignup(googleID, email, name, avatar string) (*userModels.User, error)
	Login(email, password string) (string, error)
}

type userUsecase struct {
	userRepo  userRepository.UserRepository
	secretKey string
}

func NewUserUsecase(userRepo userRepository.UserRepository, secretKey string) UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		secretKey: secretKey,
	}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("pkg::HashPassword - Error while hashing password")
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Warn().
			Str("hashedPassword", hashedPassword).
			Str("inputPassword", password).
			Err(err).
			Msg("pkg::ComparePassword - Error while comparing password")
		return false
	}
	return true
}

func (u *userUsecase) GetUsers() ([]userModels.User, error) {
	return u.userRepo.GetAll()
}

func (u *userUsecase) GetUserByID(id uuid.UUID) (userModels.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return userModels.User{}, err
	}
	if user.ID == uuid.Nil {
		return userModels.User{}, ErrNotFound
	}
	return user, nil
}

func (u *userUsecase) CreateUser(user userModels.User) (*userModels.User, error) {
	// Validasi Username
	if user.Name == "" {
		return nil, ErrUsernameValidate
	}

	// Cek apakah username sudah ada
	existingUserByUsername, err := u.userRepo.GetByUsername(user.Name)
	if err != nil {
		return nil, err
	}
	if existingUserByUsername != nil {
		return nil, ErrUsernameTaken
	}

	// Cek apakah email sudah ada
	existingUserByEmail, err := u.userRepo.GetByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	if existingUserByEmail != nil {
		return nil, ErrEmailTaken
	}

	// Hash password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return nil, err // Tangani error saat hashing password
	}
	user.Password = hashedPassword

	// Simpan pengguna baru ke database
	res, err := u.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *userUsecase) UpdateUser(user userModels.User) error {
	existingUser, err := u.userRepo.GetByID(user.ID)
	if err != nil {
		return ErrNotFound
	}

	if existingUser.Name != user.Name {
		existingUserByUsername, err := u.userRepo.GetByUsername(user.Name)
		if err != nil {
			return err
		}
		if existingUserByUsername != nil {
			return ErrUsernameTaken
		}
	}

	if existingUser.Email != user.Email {
		existingUserByEmail, err := u.userRepo.GetByEmail(user.Email)
		if err != nil {
			return err
		}
		if existingUserByEmail != nil {
			return ErrEmailTaken
		}
	}

	if user.Password != "" {
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	} else {
		user.Password = existingUser.Password
	}

	user.CreatedAt = existingUser.CreatedAt
	return u.userRepo.Update(user)
}

func (u *userUsecase) DeleteUser(id uuid.UUID) error {
	_, err := u.userRepo.GetByID(id)
	if err != nil {
		return ErrNotFound
	}
	return u.userRepo.Delete(id)
}

func (u *userUsecase) LoginOrSignup(googleID, email, name, avatar string) (*userModels.User, error) {
	var user *userModels.User
	var err error

	// Cek apakah pengguna sudah ada berdasarkan email
	user, err = u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	// Jika tidak ada berdasarkan email, cek berdasarkan Google ID
	if user == nil || user.ID == uuid.Nil {
		user, err = u.userRepo.FindGoogleId(googleID)
		if err != nil {
			return nil, err
		}
	}

	// Jika pengguna belum ada, buat pengguna baru
	if user == nil || user.ID == uuid.Nil {
		user, err = u.userRepo.Create(userModels.User{
			Name:     name,
			Email:    email,
			GoogleID: googleID,
			Avatar:   avatar,
		})
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (u *userUsecase) SearchUsers(query string) ([]userModels.User, error) {
	if query == "" {
		return nil, errors.New("search query cannot be empty")
	}
	return u.userRepo.Search(query)
}

func (u *userUsecase) Login(email, password string) (string, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil || user.ID == uuid.Nil {
		return "", ErrNotFound
	}

	if !ComparePassword(user.Password, password) {
		return "", ErrInvalidCredentials
	}

	token, err := generateJWT(user.ID, u.secretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateJWT(userID uuid.UUID, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
