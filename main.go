package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	App := fiber.New()

	App.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Hello from Fiber & mongoDB"})
	})
	App.Listen(":6000")
}
