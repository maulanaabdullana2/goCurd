package productUsecase

import (
	"errors"
	ProductModels "fiber-crud/internal/domain/product"
	ProductRepository "fiber-crud/internal/repository/product"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type ProductUsecase interface {
	FindAll() ([]ProductModels.Product, error)
	FindByID(id uuid.UUID) (ProductModels.Product, error)
	Create(product ProductModels.Product) (*ProductModels.Product, error)
	Update(product ProductModels.Product) error
	Delete(id uuid.UUID) error
}

type productUsecase struct {
	productRepository ProductRepository.ProductRepository
}

func NewProductUsecase(productRepository ProductRepository.ProductRepository) ProductUsecase {
	return &productUsecase{productRepository}
}

func (u *productUsecase) FindAll() ([]ProductModels.Product, error) {
	return u.productRepository.GetProducts()
}

func (u *productUsecase) FindByID(id uuid.UUID) (ProductModels.Product, error) {
	product, err := u.productRepository.GetProductByID(id)
	if err != nil {
		return ProductModels.Product{}, err
	}
	if product.ID == uuid.Nil {
		return ProductModels.Product{}, ErrNotFound
	}

	return u.productRepository.GetProductByID(id)
}

func (u *productUsecase) Create(product ProductModels.Product) (*ProductModels.Product, error) {
	res, err := u.productRepository.CreateProduct(product)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *productUsecase) Update(product ProductModels.Product) error {
	product, err := u.productRepository.GetProductByID(product.ID)
	if err != nil {
		return ErrNotFound
	}
	return u.productRepository.UpdateProduct(product)
}

func (u *productUsecase) Delete(id uuid.UUID) error {
	_, err := u.productRepository.GetProductByID(id)
	if err != nil {
		return ErrNotFound
	}
	return u.productRepository.DeleteProduct(id)
}
