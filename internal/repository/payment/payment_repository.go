package paymentRepository

import (
	paymentModels "fiber-crud/internal/domain/payment"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

type PaymentRepository interface {
	CreatePayment(payment *paymentModels.PaymentModels) error
	UpdatePayment(payment *paymentModels.PaymentModels) error
	GetPaymentByOrderID(orderID uuid.UUID, payment *paymentModels.PaymentModels) error
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// buat payment dengan midtrans
func (r *paymentRepository) CreatePayment(payment *paymentModels.PaymentModels) error {
	return r.db.Create(&payment).Error
}

func (r *paymentRepository) UpdatePayment(payment *paymentModels.PaymentModels) error {
	return r.db.Save(&payment).Error
}

func (r *paymentRepository) GetPaymentByOrderID(orderID uuid.UUID, payment *paymentModels.PaymentModels) error {
	return r.db.Where("order_id = ?", orderID).First(&payment).Error
}
