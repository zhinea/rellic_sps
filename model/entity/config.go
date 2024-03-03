package entity

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Mysql struct {
			DSN string `yaml:"dsn"`
		} `yaml:"mysql"`
		Redis struct {
			Addr     string `yaml:"addr"`
			Password string `yaml:"password"`
			DB       int    `yaml:"db"`
		}
	} `yaml:"database"`
}
