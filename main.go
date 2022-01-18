package main

import (
	"fiber-mongo-api/configs"
	"fiber-mongo-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// App.Get("/", func(c *fiber.Ctx) error {
	// 	return c.JSON(&fiber.Map{"data": "Hello from Fiber & mongoDB"})
	// }) g

	//run database
	configs.ConnectDB()

	//routes
	routes.UserRoute(app)
	routes.PostRoute(app)

	app.Listen(":6000")
}
