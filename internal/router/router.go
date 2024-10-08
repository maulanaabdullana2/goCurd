package router

import (
	handler "fiber-crud/internal/handler/cart"
	CommentHandler "fiber-crud/internal/handler/comment"
	paymentHandler "fiber-crud/internal/handler/payment"
	ProductHandler "fiber-crud/internal/handler/product"
	userHandler "fiber-crud/internal/handler/user"
	"fiber-crud/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, userHandler *userHandler.UserHandler) {
	app.Get("/users", userHandler.GetUsers)
	app.Get("/users/:id", userHandler.GetUserByID)
	app.Post("/users", userHandler.CreateUser)
	app.Put("/users/:id", userHandler.UpdateUser)
	app.Delete("/users/:id", userHandler.DeleteUser)
	app.Get("/search", userHandler.SearchUsers)
	app.Post("/login", userHandler.Login)
	app.Get("/auth/me", middleware.AuthMiddleware(), userHandler.CurrentUser)
	app.Get("/auth/google", userHandler.GoogleLogin)
	app.Get("/auth/google/callback", userHandler.GoogleCallback)
}

func SetupProductRoutes(app *fiber.App, productHandler *ProductHandler.ProductHandler) {
	app.Get("/products", middleware.AuthMiddleware(), productHandler.FindAll)
	app.Get("/products/:id", middleware.AuthMiddleware(), productHandler.FindByID)
	app.Post("/products", middleware.AuthMiddleware(), productHandler.Create)
	app.Put("/products/:id", middleware.AuthMiddleware(), productHandler.Update)
	app.Delete("/products/:id", middleware.AuthMiddleware(), productHandler.Delete)
	app.Get("/all-products", middleware.AuthMiddleware(), productHandler.GetAllProduct)
}

func SetupComment(app *fiber.App, commentHandler *CommentHandler.CommentHandler) {
	app.Post("/products/comments/:id", middleware.AuthMiddleware(), commentHandler.CreateCommentProductID)
	app.Get("/products/comments/:id", middleware.AuthMiddleware(), commentHandler.GetCommentsByProductid)
}

func SetupCart(app *fiber.App, cartHandler *handler.CartHandler) {
	app.Post("/carts/:id", middleware.AuthMiddleware(), cartHandler.AddItemToCart)
	app.Get("/carts", middleware.AuthMiddleware(), cartHandler.GetAllcartItems)
}

func SetupPayment(app *fiber.App, paymentHandler *paymentHandler.PaymentHandler) {
	app.Post("/payments", middleware.AuthMiddleware(), paymentHandler.CreatePayment)
	app.Post("/payment/callback", paymentHandler.UpdatePaymentStatus)
}

// func SetupUserRoutes(app *fiber.App, userHandler *userHandler.UserHandler) {
// 	app.Post("/login", userHandler.Login) // Rute login tidak memerlukan autentikasi

// 	// Rute pengguna, hanya bisa diakses oleh admin
// 	app.Get("/users", middleware.Authenticate, middleware.CheckRole("admin"), userHandler.GetUsers)
// 	app.Get("/users/:id", middleware.Authenticate, middleware.CheckRole("admin"), userHandler.GetUserByID)
// 	app.Post("/users", middleware.Authenticate, middleware.CheckRole("admin"), userHandler.CreateUser)
// 	app.Put("/users/:id", middleware.Authenticate, middleware.CheckRole("admin"), userHandler.UpdateUser)
// 	app.Delete("/users/:id", middleware.Authenticate, middleware.CheckRole("admin"), userHandler.DeleteUser)
// 	app.Get("/search", middleware.Authenticate, middleware.CheckRole("admin"), userHandler.SearchUsers)
// }

// func SetupProductRoutes(app *fiber.App, productHandler *ProductHandler.ProductHandler) {
// 	api := app.Group("/products", middleware.Authenticate) // Menambahkan middleware Authenticate ke semua rute produk

// 	// Rute produk, bisa diakses oleh admin dan user
// 	api.Get("/", middleware.CheckRole("admin", "user"), productHandler.FindAll)
// 	api.Get("/:id", middleware.CheckRole("admin", "user"), productHandler.FindByID)
// 	api.Post("/", middleware.CheckRole("admin"), productHandler.Create)
// 	api.Put("/:id", middleware.CheckRole("admin"), productHandler.Update)
// 	api.Delete("/:id", middleware.CheckRole("admin"), productHandler.Delete)
// }
