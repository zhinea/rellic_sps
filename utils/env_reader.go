package utils

import (
	"github.com/zhinea/sps/model/entity"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
)

func EnvReader(configPath string) *entity.Config {

	log.Println("Load config from =>", configPath)

	f, err := os.Open(configPath)
	defer f.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var cfg entity.Config

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalln(err)
	}

	return &cfg
}

func GetEnvPath() string {
	dir, err := os.Getwd()

	if err != nil {
		log.Println(err)
		panic(err)
	}

	return path.Join(dir, "config.yml")
}
