package handler

import (
	usecase "fiber-crud/internal/usecase/cart"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CartHandler struct {
	cartUsecase usecase.CartUsecase
}

func NewCartHandler(cartUsecase usecase.CartUsecase) *CartHandler {
	return &CartHandler{
		cartUsecase: cartUsecase,
	}
}

func (h *CartHandler) AddItemToCart(c *fiber.Ctx) error {
	productIDStr := c.Params("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID format"})
	}

	var request struct {
		Quantity int `json:"quantity"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if request.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Quantity must be greater than zero"})
	}

	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	err = h.cartUsecase.AddItemToCart(userID, productID, request.Quantity)
	if err != nil {
		if err == usecase.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Item added to cart successfully"})
}
