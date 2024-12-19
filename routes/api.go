package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radityajay/go-url-shortener/controllers"
)

func RouteApiRegister(app *fiber.App) {
	// Routing in here
	app.Get("/", controllers.Main)

	// Routing shortener
	app.Post("/shortener", controllers.CreateShortURL)

	route := app.Group("/s")
	route.Get("/:short_code", controllers.GetShortURL)

}
