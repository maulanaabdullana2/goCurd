package userHandler

import (
	"context"
	"encoding/json"
	"net/http"

	userModels "fiber-crud/internal/domain/user"
	Userusecase "fiber-crud/internal/usecase/user"
	"fiber-crud/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type UserHandler struct {
	userUsecase Userusecase.UserUsecase
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(usecase Userusecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: usecase}
}

// GetUsers handles requests to get all users
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.userUsecase.GetUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

func (h *UserHandler) CurrentUser(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := h.userUsecase.GetCurrentUser(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user"})
	}

	return c.JSON(user)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID format"})
	}

	user, err := h.userUsecase.GetUserByID(id)
	if err == Userusecase.ErrNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user userModels.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.userUsecase.CreateUser(user)
	if err == Userusecase.ErrUsernameTaken || err == Userusecase.ErrEmailTaken {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID format"})
	}

	var user userModels.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	user.ID = id

	err = h.userUsecase.UpdateUser(user)
	if err == Userusecase.ErrNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	} else if err == Userusecase.ErrUsernameTaken || err == Userusecase.ErrEmailTaken {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID format"})
	}

	err = h.userUsecase.DeleteUser(id)
	if err == Userusecase.ErrNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) SearchUsers(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Search query cannot be empty"})
	}

	users, err := h.userUsecase.SearchUsers(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	token, err := h.userUsecase.Login(credentials.Email, credentials.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": token})
}

func (h *UserHandler) GoogleLogin(c *fiber.Ctx) error {
	url := utils.GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

func (h *UserHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	token, err := utils.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	client := utils.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user info"})
	}

	var googleUser struct {
		GoogleID string `json:"sub"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Avatar   string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.userUsecase.LoginOrSignup(
		googleUser.GoogleID,
		googleUser.Email,
		googleUser.Name,
		googleUser.Avatar,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	tokenString, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": tokenString, "user": user})
}
