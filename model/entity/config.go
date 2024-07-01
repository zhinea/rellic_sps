package entity

type Config struct {
	Container struct {
		ServerID  string `yaml:"server_id"`
		ServerUrl string `yaml:"server_url"`
	} `yaml:"container"`

	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`

		Domain     string `yaml:"domain"`
		SystemPath string `yaml:"system_path"`
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
