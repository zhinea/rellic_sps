package routes

import "github.com/gofiber/fiber/v2"

func RouteInit(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "What are you doing here?",
		})
	})
}
