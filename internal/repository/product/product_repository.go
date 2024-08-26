package ProductRepository

import (
	ProductModels "fiber-crud/internal/domain/product"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProducts(userID uuid.UUID) ([]ProductModels.Product, error)
	GetProductByID(id uuid.UUID, userID uuid.UUID) (ProductModels.Product, error)
	CreateProduct(product *ProductModels.Product) (*ProductModels.Product, error)
	UpdateProduct(product *ProductModels.Product) error
	DeleteProduct(id uuid.UUID, userID uuid.UUID) error
	DecreaseStock(productID uuid.UUID, quantity int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetProducts(userID uuid.UUID) ([]ProductModels.Product, error) {
	var products []ProductModels.Product
	if err := r.db.Preload("Comments").Where("user_id = ?", userID).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) GetProductByID(id uuid.UUID, userID uuid.UUID) (ProductModels.Product, error) {
	var product ProductModels.Product
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ProductModels.Product{}, nil
		}
		return ProductModels.Product{}, err
	}
	return product, nil
}

func (r *productRepository) CreateProduct(product *ProductModels.Product) (*ProductModels.Product, error) {
	if err := r.db.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) UpdateProduct(product *ProductModels.Product) error {
	if err := r.db.Save(product).Error; err != nil {
		return err
	}
	return nil
}

func (r *productRepository) DeleteProduct(id uuid.UUID, userID uuid.UUID) error {
	if err := r.db.Delete(&ProductModels.Product{}, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return err
	}
	return nil
}

func (r *productRepository) DecreaseStock(productID uuid.UUID, quantity int) error {

	result := r.db.Model(&ProductModels.Product{}).
		Where("id = ? AND stock >= ?", productID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
