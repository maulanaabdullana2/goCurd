package paymentHandler

import (
	"net/http"

	paymentUsecase "fiber-crud/internal/usecase/payment"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	usecase paymentUsecase.PaymentUsecase
}

func NewPaymentHandler(usecase paymentUsecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		usecase: usecase,
	}
}

func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	redirectURL, err := h.usecase.CreatePaymentMidtrans(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"redirect_url": redirectURL,
	})
}

func (h *PaymentHandler) UpdatePaymentStatus(c *fiber.Ctx) error {
	var callbackData struct {
		OrderID string `json:"order_id"`
		Status  string `json:"transaction_status"`
	}

	if err := c.BodyParser(&callbackData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if len(callbackData.OrderID) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order ID cannot be empty",
		})
	}

	orderID, err := uuid.Parse(callbackData.OrderID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID format",
		})
	}

	err = h.usecase.UpdatePaymentstatus(orderID, callbackData.Status)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Payment status updated successfully",
	})
}
