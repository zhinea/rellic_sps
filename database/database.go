package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/zhinea/sps/model/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB
var Redis *redis.Client
var Ctx = context.Background()

func InitDatabase(cfg *entity.Config) {
	var err error

	DB, err = gorm.Open(mysql.Open(cfg.Database.Mysql.DSN), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	log.Println("Connection Opened to Database")

	Redis = redis.NewClient(&redis.Options{
		Addr:     cfg.Database.Redis.Addr,
		Password: cfg.Database.Redis.Password,
		DB:       cfg.Database.Redis.DB,
	})

	//defer Redis.Close()

	//defer func(r *redis.Client) {
	//	if err := r.Close(); err != nil {
	//		log.Fatal(err)
	//	}
	//}(Redis)

	// Perform basic diagnostic to check if the connection is working
	// Expected result > ping: PONG
	// If Redis is not running, error case is taken instead
	if errRedis := Redis.Ping(Ctx).Err(); errRedis != nil {
		log.Fatalln("Redis connection was refused")
	} else {
		log.Println("Redis connection was successful")
	}
}
