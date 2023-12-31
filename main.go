package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zhinea/sps/routes"
)

func main() {
	app := fiber.New()

	routes.RouteInit(app)

	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
