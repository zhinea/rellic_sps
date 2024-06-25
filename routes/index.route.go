package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zhinea/sps/controllers/gtagcontroller"
)

func RouteInit(app *fiber.App) {

	app.Get("/gtag/js", gtagcontroller.GetScripts)
	app.All("/:any", gtagcontroller.HandleTrackData)

}
