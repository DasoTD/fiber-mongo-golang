package main

import (
	"fiber-mongo-api/configs"

	"github.com/gofiber/fiber/v2"
)

func main() {
	App := fiber.New()

	// App.Get("/", func(c *fiber.Ctx) error {
	// 	return c.JSON(&fiber.Map{"data": "Hello from Fiber & mongoDB"})
	// }) g

	//run database
	configs.ConnectDB()

	App.Listen(":6000")
}
