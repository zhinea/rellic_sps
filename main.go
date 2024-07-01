package main

import (
	"database/sql"
	"flag"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/robfig/cron/v3"
	"github.com/zhinea/sps/cronjob"
	"github.com/zhinea/sps/database"
	"github.com/zhinea/sps/handler"
	"github.com/zhinea/sps/model/entity"
	"github.com/zhinea/sps/routes"
	"github.com/zhinea/sps/utils"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	cfg     entity.Config
	once    sync.Once
	sqlDB   *sql.DB
	initErr error
)

func initialize() {
	cfgFilename := flag.String("config", utils.GetEnvPath(), "Config file path.")
	flag.Parse()

	cfg = *utils.EnvReader(*cfgFilename)

	// Initial database
	database.InitDatabase(&cfg)

	// Retrieve sql.DB from GORM DB
	sqlDB, initErr = database.DB.DB()
	if initErr != nil {
		log.Fatalln(initErr)
	}
}

func main() {
	// Ensure initialization runs only once
	once.Do(initialize)

	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		Prefork:      true,
		ServerHeader: "Proxy Server by rellic.app",
	})

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	sysRoute := app.Group(cfg.Server.SystemPath)
	routes.SysRouteInit(sysRoute)

	// Initial middleware
	app.Use(adaptor.HTTPMiddleware(handler.AppMiddleware))

	// Initial routes
	routes.RouteInit(app)

	// Close database connections on shutdown
	defer sqlDB.Close()
	defer database.Redis.Close()

	// Schedule cron job only in main process
	if !fiber.IsChild() {
		log.Println("Handling cronjob", syscall.Getpid())
		jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
		scheduler := cron.New(cron.WithLocation(jakartaTime))

		defer scheduler.Stop()

		scheduler.AddFunc("*/5 * * * *", cronjob.BillingSchedule)

		go scheduler.Start()
	}

	err := app.Listen(cfg.Server.Host + ":" + cfg.Server.Port)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// Trap SIGINT to trigger shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
