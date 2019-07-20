package config

import (
	"github.com/joeshaw/envdecode"
	"github.com/subosito/gotenv"
)

type Config struct {
	Env  string `env:"ENV"`
	Port uint16 `env:"PORT,default=8081"`

	Database struct {
		Host     string `env:"DATABASE_HOST,default=localhost"`
		Name     string `env:"DATABASE_NAME,required"`
		Username string `env:"DATABASE_USERNAME,required"`
		Password string `env:"DATABASE_PASSWORD,required"`
		Pool     int    `env:"DATABASE_POOL,default=5000"`
	}

	Redis struct {
		Address string `env:"REDIS_HOST,required"`
	}
}

func NewConfig() *Config {
	var cfg Config
	gotenv.Load(".env")
	err := envdecode.Decode(&cfg)
	check(err)
	return &cfg
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func RunServer() {
	//TODO Listen & Serve
}
