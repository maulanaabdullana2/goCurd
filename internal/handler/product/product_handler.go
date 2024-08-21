package ProductHandler

import (
	ProductModels "fiber-crud/internal/domain/product"
	productUsecase "fiber-crud/internal/usecase/product"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productUsecase productUsecase.ProductUsecase
}

func NewProductHandler(productUsecase productUsecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
	}
}

func (h *ProductHandler) FindAll(c *fiber.Ctx) error {
	products, err := h.productUsecase.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(products)
}

func (h *ProductHandler) FindByID(c *fiber.Ctx) error {
	idstr := c.Params("id")
	id, err := uuid.Parse(idstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	product, err := h.productUsecase.FindByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})

	}

	return c.JSON(product)
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var product ProductModels.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.productUsecase.Create(product)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(res)
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	var product ProductModels.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	idstr := c.Params("id")
	id, err := uuid.Parse(idstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	product.ID = id
	err = h.productUsecase.Update(product)

	if err == productUsecase.ErrNotFound {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(product)
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	idstr := c.Params("id")
	id, err := uuid.Parse(idstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	err = h.productUsecase.Delete(id)
	if err == productUsecase.ErrNotFound {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
