package main

import (
	handler "fiber-crud/internal/handler/cart"
	commentHandler "fiber-crud/internal/handler/comment"
	ProductHandler "fiber-crud/internal/handler/product"
	UserHandel "fiber-crud/internal/handler/user"
	user "fiber-crud/internal/repository"
	CartRepository "fiber-crud/internal/repository/cart"
	repository "fiber-crud/internal/repository/comment"
	ProductRepository "fiber-crud/internal/repository/product"
	"fiber-crud/internal/router"
	usecase "fiber-crud/internal/usecase/cart"
	commentUsecase "fiber-crud/internal/usecase/comment"
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
	userUsecase := Userusecase.NewUserUsecase(userRepo)
	userHandler := UserHandel.NewUserHandler(userUsecase)

	productRepo := ProductRepository.NewProductRepository(db)
	productUsecase := productUsecase.NewProductUsecase(productRepo)
	productHandler := ProductHandler.NewProductHandler(productUsecase)

	commentRepo := repository.NewCommentRepository(db)
	commentUsecase := commentUsecase.NewCommentUsecase(commentRepo)
	commentHandler := commentHandler.NewCommentHandler(commentUsecase)

	cartRepo := CartRepository.NewCartRepository(db)
	cartUsecase := usecase.NewCartUsecase(cartRepo, productRepo)
	cartHandler := handler.NewCartHandler(cartUsecase)

	app := fiber.New()

	router.SetupUserRoutes(app, userHandler)
	router.SetupProductRoutes(app, productHandler)
	router.SetupComment(app, commentHandler)
	router.SetupCart(app, cartHandler)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Route Not Found",
		})
	})

	app.Listen(":3000")
}
