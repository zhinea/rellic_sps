package utils

import (
	"log"
	"os"
	"path"
)

func EnvReader() (*os.File, error) {
	dir, err := os.Getwd()

	if err != nil {
		log.Println(err)
		panic(err)
	}
	configPath := path.Join(dir, "config.yml")

	log.Println("Load config from =>", configPath)

	return os.Open(configPath)
}
