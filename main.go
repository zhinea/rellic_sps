package main

import (
	"flag"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/zhinea/sps/database"
	"github.com/zhinea/sps/handler"
	"github.com/zhinea/sps/routes"
	"github.com/zhinea/sps/utils"
	"log"
)

func main() {
	cfgFilename := flag.String("config", utils.GetEnvPath(), "Config file path.")

	cfg := *utils.EnvReader(*cfgFilename)

	// initial database
	database.InitDatabase(&cfg)

	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		Prefork:      true,
		ServerHeader: "Proxy Server by rellic.app",
	})

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	sysRoute := app.Group("/c0n9wb-sys")

	routes.SysRouteInit(sysRoute)

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

	err = app.Listen(":" + cfg.Server.Port)
	if err != nil {
		return
	}
}
