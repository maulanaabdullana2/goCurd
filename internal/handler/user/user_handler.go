package userHandler

import (
	userModels "fiber-crud/internal/domain/user"
	Userusecase "fiber-crud/internal/usecase/user"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	userUsecase Userusecase.UserUsecase
}

// NewUserHandler membuat instance baru UserHandler
func NewUserHandler(usecase Userusecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: usecase}
}

// GetUsers menangani request untuk mengambil semua pengguna
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.userUsecase.GetUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// GetUserByID menangani request untuk mengambil pengguna berdasarkan ID
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid UUID format"})
	}

	user, err := h.userUsecase.GetUserByID(id)
	if err == Userusecase.ErrNotFound {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// CreateUser menangani request untuk membuat pengguna baru
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user userModels.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.userUsecase.CreateUser(user)
	if err == Userusecase.ErrUsernameTaken {
		return c.Status(409).JSON(fiber.Map{"err": err.Error()})
	} else if err == Userusecase.ErrEmailTaken {
		return c.Status(409).JSON(fiber.Map{"error": err.Error()})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(res)
}

// UpdateUser menangani request untuk memperbarui pengguna
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid UUID format"})
	}

	var user userModels.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	user.ID = id

	err = h.userUsecase.UpdateUser(user)
	if err == Userusecase.ErrNotFound {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	} else if err == Userusecase.ErrUsernameTaken {
		return c.Status(409).JSON(fiber.Map{"error": "Username already taken"})
	} else if err == Userusecase.ErrEmailTaken {
		return c.Status(409).JSON(fiber.Map{"error": "Email already taken"})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

// DeleteUser menangani request untuk menghapus pengguna
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid UUID format"})
	}

	err = h.userUsecase.DeleteUser(id)
	if err == Userusecase.ErrNotFound {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}

// SearchUsers menangani request untuk mencari pengguna berdasarkan query
func (h *UserHandler) SearchUsers(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Search query cannot be empty"})
	}

	users, err := h.userUsecase.SearchUsers(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// Login menangani request login pengguna
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	token, err := h.userUsecase.Login(credentials.Email, credentials.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"token":   token,
			"message": "Login successful",
		},
	})
}
