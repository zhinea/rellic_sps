package main

import (
	"context"
	database "github.com/zhinea/sps/utils"
	"log"
	"testing"
)

func TestConnect(t *testing.T) {

	db := database.GetConnection()

	if db == nil {
		t.Errorf("db is nil")
	}

	err := db.Ping()
	if err != nil {
		t.Errorf("db ping error: %s", err.Error())
	}

	log.Println("db ping success")
}

func TestRedisConnect(t *testing.T) {

	db := database.GetRedis()

	if db == nil {
		t.Errorf("db is nil")
	}

	err := db.Ping(context.Background()).Err()
	if err != nil {
		t.Errorf("db ping error: %s", err.Error())
	}

	log.Println("db ping success")
}

func TestExecWithCache(t *testing.T) {

	db := database.GetConnection()
	ctx := context.Background()

	if db == nil {
		t.Errorf("db is nil")
	}

	res, err := database.ExecWithCache(ctx, "SELECT * FROM users")
	if err != nil {
		t.Errorf("db exec error: %s", err.Error())
	}

	log.Println("db exec success", res)
}
