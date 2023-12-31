package database

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func GetConnection() *sql.DB {
	dbConnection := "root:admin@tcp(localhost:3306)/rellic_iofi"
	db, err := sql.Open("mysql", dbConnection)

	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func GetRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0, // use default DB
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return client
}

func ExecWithCache(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	client := GetRedis()

	cacheArgs := args

	if args == nil {
		cacheArgs = []interface{}{""}
	}

	cacheKey := "query:" + query + ":args:" + cacheArgs[0].(string)

	log.Println(cacheKey)

	log.Println("query", query)

	//check cache
	cacheRes, cacheErr := client.Get(ctx, cacheKey).Result()
	if cacheErr == nil && cacheRes != "" {
		log.Println("HIT Cache")
		return cacheRes, nil
	}

	db := GetConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	log.Println("MISS Cache")

	res, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	//cache before send
	client.Set(ctx, cacheKey, res, -1)

	return res, nil
}
