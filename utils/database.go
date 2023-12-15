package database

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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

func GetCacheConnection() {
	client := redis
}

func ExecWithCache(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	db := GetConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	res, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
