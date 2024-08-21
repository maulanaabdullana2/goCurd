package productUsecase

import (
	"errors"
	ProductModels "fiber-crud/internal/domain/product"
	ProductRepository "fiber-crud/internal/repository/product"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("product not found")

type ProductUsecase interface {
	GetProducts(userID uuid.UUID) ([]ProductModels.Product, error)
	GetProductByID(id uuid.UUID, userID uuid.UUID) (ProductModels.Product, error)
	CreateProduct(product *ProductModels.Product) (*ProductModels.Product, error)
	UpdateProduct(product *ProductModels.Product, userID uuid.UUID) error
	DeleteProduct(id uuid.UUID, userID uuid.UUID) error
}

type productUsecase struct {
	productRepo ProductRepository.ProductRepository
}

func NewProductUsecase(repo ProductRepository.ProductRepository) ProductUsecase {
	return &productUsecase{productRepo: repo}
}

func (u *productUsecase) GetProducts(userID uuid.UUID) ([]ProductModels.Product, error) {
	return u.productRepo.GetProducts(userID)
}

func (u *productUsecase) GetProductByID(id uuid.UUID, userID uuid.UUID) (ProductModels.Product, error) {
	product, err := u.productRepo.GetProductByID(id, userID)
	if err != nil || product.ID == uuid.Nil {
		return ProductModels.Product{}, ErrNotFound
	}
	return product, nil
}

func (u *productUsecase) CreateProduct(product *ProductModels.Product) (*ProductModels.Product, error) {
	return u.productRepo.CreateProduct(product)
}

func (u *productUsecase) UpdateProduct(product *ProductModels.Product, userID uuid.UUID) error {
	existingProduct, err := u.productRepo.GetProductByID(product.ID, userID)
	if err != nil {
		return err
	}
	if existingProduct.ID == uuid.Nil {
		return ErrNotFound
	}
	product.UserID = userID
	return u.productRepo.UpdateProduct(product)
}

func (u *productUsecase) DeleteProduct(id uuid.UUID, userID uuid.UUID) error {
	existingProduct, err := u.productRepo.GetProductByID(id, userID)
	if err != nil {
		return err
	}
	if existingProduct.ID == uuid.Nil {
		return ErrNotFound
	}
	return u.productRepo.DeleteProduct(id, userID)
}
