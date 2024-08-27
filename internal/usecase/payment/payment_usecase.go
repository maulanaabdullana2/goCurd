package paymentUsecase

import (
	"errors"
	paymentModels "fiber-crud/internal/domain/payment"
	cartRepository "fiber-crud/internal/repository/cart"
	paymentRepository "fiber-crud/internal/repository/payment"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/veritrans/go-midtrans"
)

type PaymentUsecase interface {
	UpdatePaymentstatus(orderID uuid.UUID, status string) error
	CreatePaymentMidtrans(userID uuid.UUID) (string, error)
}

type paymentUsecase struct {
	paymentRepo paymentRepository.PaymentRepository
	cartRepo    cartRepository.CartRepository
	midtrans    midtrans.Client
}

func NewPaymentUsecase(paymentRepo paymentRepository.PaymentRepository, cartRepo cartRepository.CartRepository) PaymentUsecase {
	midtransServerKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if midtransServerKey == "" {
		panic("Midtrans server key not set in environment variables")
	}

	midtransClient := midtrans.NewClient()
	midtransClient.ServerKey = midtransServerKey
	midtransClient.ClientKey = os.Getenv("MIDTRANS_CLIENT_KEY")
	midtransClient.APIEnvType = midtrans.Sandbox

	return &paymentUsecase{
		paymentRepo: paymentRepo,
		cartRepo:    cartRepo,
		midtrans:    midtransClient,
	}
}

func (p *paymentUsecase) CreatePaymentMidtrans(userID uuid.UUID) (string, error) {
	carts, err := p.cartRepo.GetAllcartItems(userID)
	if err != nil {
		return "", err
	}

	if len(carts) == 0 {
		return "", errors.New("no items in cart")
	}

	var total int
	for _, cart := range carts {
		total += int(cart.Product.Price) * cart.Quantity
	}

	orderID := uuid.New().String()

	payment := &paymentModels.PaymentModels{
		ID:      uuid.New(),
		OrderID: orderID,
		UserID:  userID,
		Status:  "pending",
		Amount:  total,
	}

	err = p.paymentRepo.CreatePayment(payment)
	if err != nil {
		return "", err
	}

	params := midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(total),
		},
	}

	snapGateway := midtrans.SnapGateway{Client: p.midtrans}
	snapResp, err := snapGateway.GetToken(&params)
	if err != nil {
		return "", err
	}

	return snapResp.RedirectURL, nil
}

func (p *paymentUsecase) UpdatePaymentstatus(orderID uuid.UUID, status string) error {
	if orderID == uuid.Nil {
		return errors.New("orderID cannot be empty")
	}

	payment := &paymentModels.PaymentModels{}
	err := p.paymentRepo.GetPaymentByOrderID(orderID, payment)
	if err != nil {
		return fmt.Errorf("failed to fetch payment: %v", err)
	}

	payment.Status = status

	err = p.paymentRepo.UpdatePayment(payment)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %v", err)
	}

	return nil
}
