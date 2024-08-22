package main

import (
	ProductHandler "fiber-crud/internal/handler/product"
	UserHandel "fiber-crud/internal/handler/user"
	user "fiber-crud/internal/repository"
	ProductRepository "fiber-crud/internal/repository/product"
	"fiber-crud/internal/router"
	productUsecase "fiber-crud/internal/usecase/product"
	Userusecase "fiber-crud/internal/usecase/user"
	db "fiber-crud/package"
	"fiber-crud/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {

	utils.InitOAuth2()
	utils.InitCloudinary()
	db := db.InitDB()

	userRepo := user.NewUserRepository(db)
	secretKey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	userUsecase := Userusecase.NewUserUsecase(userRepo, secretKey)
	userHandler := UserHandel.NewUserHandler(userUsecase)

	productRepo := ProductRepository.NewProductRepository(db)
	productUsecase := productUsecase.NewProductUsecase(productRepo)
	productHandler := ProductHandler.NewProductHandler(productUsecase)

	app := fiber.New()

	router.SetupUserRoutes(app, userHandler)
	router.SetupProductRoutes(app, productHandler)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Route Not Found",
		})
	})

	app.Listen(":3000")
}
