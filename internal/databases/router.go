package databases

import "github.com/gofiber/fiber/v2"

func AddTodoRoutes(app *fiber.App, controller *DatabasesController) {
	databases := app.Group("/databases")

	databases.Post("/", controller.create)
}