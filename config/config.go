package config

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/faris-arifiansyah/fws-rsvp/delivery"
	"github.com/faris-arifiansyah/fws-rsvp/handler"
	"github.com/faris-arifiansyah/fws-rsvp/repository"
	"github.com/faris-arifiansyah/fws-rsvp/usecase"
	"github.com/faris-arifiansyah/mgoi"
	"github.com/go-redis/redis"
	"github.com/joeshaw/envdecode"
	"github.com/rs/cors"
	"github.com/subosito/gotenv"
)

type Config struct {
	Env  string `env:"ENV"`
	Port uint16 `env:"PORT,default=8082"`

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

type RedisOption struct {
	Address      string
	PingTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MaxRetries   int
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

func NewRedis(opt RedisOption) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         opt.Address,
		DialTimeout:  opt.PingTimeout,
		ReadTimeout:  opt.ReadTimeout,
		WriteTimeout: opt.WriteTimeout,
	})

	_, err := client.Ping().Result()

	return client, err
}

func RunServer() {
	cfg := NewConfig()

	//dependencies
	db, err := NewMongoDB(cfg)
	check(err)

	redisOpt := RedisOption{
		Address:      cfg.Redis.Address,
		PingTimeout:  time.Duration(1 * time.Second),
		ReadTimeout:  time.Duration(1 * time.Second),
		WriteTimeout: time.Duration(1 * time.Second),
		MaxRetries:   3,
	}

	redis, err := NewRedis(redisOpt)
	check(err)

	rsvpRepo := repository.NewMongoRsvp(db)
	uc := usecase.NewRsvpUsecase(&usecase.AccessProvider{
		RsvpRepo: rsvpRepo,
	})

	rsvpHandler := delivery.NewRsvpHandler(uc, redis)
	h, err := handler.NewHandler(&rsvpHandler)
	check(err)

	co := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "PUT", "HEAD", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		MaxAge:         86400,
	})

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      co.Handler(h),
		ReadTimeout:  310 * time.Second,
		WriteTimeout: 310 * time.Second,
	}

	log.Printf("RSVP is available at %s\n", s.Addr)
	if serr := s.ListenAndServe(); serr != http.ErrServerClosed {
		log.Fatal(serr)
	}
}
