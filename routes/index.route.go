package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zhinea/sps/controllers/gtagcontroller"
)

func RouteInit(app *fiber.App) {

	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status": 200,
		})
	})

	app.Get("/gtag/js", gtagcontroller.GetScripts)
	app.All("/:any", gtagcontroller.HandleTrackData)

}
