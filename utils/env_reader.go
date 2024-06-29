package utils

import (
	"github.com/zhinea/sps/model/entity"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
)

var Cfg entity.Config

func EnvReader(configPath string) *entity.Config {

	log.Println("Load config from =>", configPath)

	f, err := os.Open(configPath)
	defer f.Close()

	if err != nil {
		log.Fatalln(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Cfg)
	if err != nil {
		log.Fatalln(err)
	}

	return &Cfg
}

func GetEnvPath() string {
	dir, err := os.Getwd()

	if err != nil {
		log.Println(err)
		panic(err)
	}

	return path.Join(dir, "config.yml")
}
