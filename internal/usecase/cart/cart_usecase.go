package usecase

import (
	"errors"
	cartModels "fiber-crud/internal/domain/cart"
	CartRepository "fiber-crud/internal/repository/cart"
	ProductRepository "fiber-crud/internal/repository/product"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("product not found")

type cartUsecase struct {
	cartRepository    CartRepository.CartRepository
	productRepository ProductRepository.ProductRepository
}

type CartUsecase interface {
	AddItemToCart(userID uuid.UUID, productID uuid.UUID, quantity int) error
	GetAllcartItems(userID uuid.UUID) ([]cartModels.CartModels, error)
}

func NewCartUsecase(cartRepo CartRepository.CartRepository, productRepo ProductRepository.ProductRepository) CartUsecase {
	return &cartUsecase{cartRepo, productRepo}
}

func (u *cartUsecase) AddItemToCart(userID uuid.UUID, productID uuid.UUID, quantity int) error {
	// Fetch product to check if it exists
	product, err := u.productRepository.GetAllProductsByid(productID)
	if err != nil {
		return err
	}
	if len(product) == 0 {
		return ErrNotFound
	}

	cartItem, err := u.cartRepository.GetCartItemByProductID(userID, productID)
	if err != nil && err != CartRepository.ErrNotFound {
		return err
	}

	// If cart item exists, update its quantity
	if cartItem.ID != uuid.Nil {
		cartItem.Quantity += quantity
		if err := u.cartRepository.UpdateCartItem(cartItem); err != nil {
			return err
		}
	} else {
		// If cart item does not exist, create a new entry
		cartItem = cartModels.CartModels{
			UserID:    userID,
			ProductID: productID,
			Quantity:  quantity,
		}
		if err := u.cartRepository.AddItemToCart(userID, productID, quantity); err != nil {
			return err
		}
	}

	// Decrease stock after updating or adding the cart item
	if err := u.productRepository.DecreaseStock(productID, quantity); err != nil {
		return err
	}

	return nil
}

func (u *cartUsecase) GetAllcartItems(userID uuid.UUID) ([]cartModels.CartModels, error) {
	return u.cartRepository.GetAllcartItems(userID)
}
