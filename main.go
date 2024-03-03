package main

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/zhinea/sps/database"
	"github.com/zhinea/sps/handler"
	"github.com/zhinea/sps/routes"
	"log"
)

func main() {
	// initial database
	database.InitDatabase()

	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		Prefork:      true,
		ServerHeader: "Proxy Server by rellic.app",
	})

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// initial middleware
	app.Use(adaptor.HTTPMiddleware(handler.AppMiddleware))

	// initial routes
	routes.RouteInit(app)

	//defer database.DB.Close()
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatalln(err)
	}

	defer sqlDB.Close()
	defer database.Redis.Close()

	err = app.Listen(":3000")
	if err != nil {
		return
	}
}
