package CartRepository

import (
	cartModels "fiber-crud/internal/domain/cart"

	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("cart item not found")

type CartRepository interface {
	GetCartItemByProductID(userID uuid.UUID, productID uuid.UUID) (cartModels.CartModels, error)
	AddItemToCart(userID uuid.UUID, productID uuid.UUID, quantity int) error
	UpdateCartItem(cartItem cartModels.CartModels) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db}
}

func (r *cartRepository) GetCartItemByProductID(userID uuid.UUID, productID uuid.UUID) (cartModels.CartModels, error) {
	var cartItem cartModels.CartModels
	if err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return cartItem, ErrNotFound
		}
		return cartItem, err
	}
	return cartItem, nil
}

func (r *cartRepository) AddItemToCart(userID uuid.UUID, productID uuid.UUID, quantity int) error {
	var cartItem cartModels.CartModels
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&cartItem).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// If item does not exist, create it
		cartItem = cartModels.CartModels{
			UserID:    userID,
			ProductID: productID,
			Quantity:  quantity,
		}
		return r.db.Create(&cartItem).Error
	} else {
		// If item exists, update its quantity
		cartItem.Quantity += quantity
		return r.db.Save(&cartItem).Error
	}
}

func (r *cartRepository) UpdateCartItem(cartItem cartModels.CartModels) error {
	return r.db.Save(&cartItem).Error
}
