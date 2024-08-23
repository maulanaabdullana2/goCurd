package commentHandler

import (
	CommentModels "fiber-crud/internal/domain/comment"
	commentUsecase "fiber-crud/internal/usecase/comment"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CommentHandler struct {
	commentUsecase commentUsecase.CommentUsecase
}

func NewCommentHandler(usecase commentUsecase.CommentUsecase) *CommentHandler {
	return &CommentHandler{
		commentUsecase: usecase,
	}
}

func (h *CommentHandler) CreateCommentProductID(c *fiber.Ctx) error {
	// Parse the product ID from the URL parameters
	idstr := c.Params("id")
	id, err := uuid.Parse(idstr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID format"})
	}

	// Extract the user ID from the context
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Define a struct to parse the request body
	type CommentRequest struct {
		Content string `json:"content"`
	}

	// Parse the request body
	var requestBody CommentRequest
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Create the comment model
	comment := &CommentModels.Comment{
		ProductID: id,
		UserID:    userID,
		Content:   requestBody.Content,
	}

	// Use the comment usecase to create the comment
	err = h.commentUsecase.CreateComment(comment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return success status with comment data
	return c.JSON(fiber.Map{
		"status":  "success",
		"comment": comment,
	})
}

func (h *CommentHandler) GetCommentsByProductid(c *fiber.Ctx) error {
	idstr := c.Params("id")
	id, err := uuid.Parse(idstr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID format"})
	}

	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	comments, err := h.commentUsecase.Getcommentproductid(id, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(comments)
}
