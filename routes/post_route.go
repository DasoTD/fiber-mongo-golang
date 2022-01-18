package routes

import (
	"fiber-mongo-api/postController"

	"github.com/gofiber/fiber/v2"
)

// type postResponse struct {
// 	Status  int        `json:"status"`
// 	Message string     `json:"message"`
// 	Data    *fiber.Map `json:"data"`
// }

func PostRoute(app *fiber.App) {
	//All routes related to posts comes here
	app.Post("/post", postController.CreatePost)

	app.Get("/post/:postId", postController.GetPost)
	app.Put("/post/:postId", postController.EditPost)
	app.Delete("/post/:postId", postController.DeletePost)
	app.Get("/posts", postController.GetAllPost)
}
