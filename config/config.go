package config

import (
	"fmt"
	"time"

	"github.com/faris-arifiansyah/fws-rsvp/repository"
	"github.com/faris-arifiansyah/fws-rsvp/usecase"
	"github.com/faris-arifiansyah/mgoi"
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

func NewMongoDB(cfg *Config) (mgoi.DatabaseManager, error) {
	if cfg.Env == "test" {
		return &mgoi.Database{}, nil
	}

	fmt.Printf("Connecting to mongodb://[USERNAME]:[PASSWORD]@%s/%s --authenticationDatabase\n", cfg.Database.Host, cfg.Database.Name)

	dialInfo := &mgoi.DialInfo{
		Addrs:    []string{cfg.Database.Host},
		Username: cfg.Database.Username,
		Password: cfg.Database.Password,
		Timeout:  2 * time.Second,
	}

	dialer := mgoi.NewDialer()
	session, err := dialer.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	db := session.DB(cfg.Database.Name)

	return db, nil
}

func RunServer() {
	cfg := NewConfig()

	//dependencies
	db, err := NewMongoDB(cfg)
	check(err)

	rsvpRepo := repository.NewMongoRsvp(db)
	uc := usecase.NewRsvpUsecase(&usecase.AccessProvider{
		RsvpRepo: rsvpRepo,
	})

	fmt.Println("UC : ", uc)
	//TODO Listen & Serve
}
