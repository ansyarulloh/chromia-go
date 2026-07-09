package main

import (
	"log"

	"chromia-api/controllers"
	"chromia-api/database"
	"chromia-api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	database.ConnectDB()

	app := fiber.New()
	app.Use(logger.New())

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Server Chromia Go Fiber jalan bro sat set!",
		})
	})

	api := app.Group("/api")

    // 1. Rute Publik (Tanpa Satpam)
    auth := api.Group("/auth")
    auth.Post("/login", controllers.Login)
    auth.Post("/register", controllers.Register) // <-- TAMBAHIN BARIS INI BRO

    // 2. Rute Private (Wajib bawa Token JWT)
    protected := api.Group("/", middleware.Protected())

    user := protected.Group("/user")
    user.Get("/status", controllers.GetUserStatus)

    protected.Post("/scan", controllers.ScanColor)

    payment := protected.Group("/payment")
    payment.Post("/subscribe", controllers.SubscribePremium)

    log.Fatal(app.Listen(":3000"))
}