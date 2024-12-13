package main

import (
	"emailregister/services/register/internal/db"
	"emailregister/services/register/internal/redis"
	"emailregister/services/register/internal/router"
	"emailregister/services/register/internal/sender"
	"emailregister/services/register/internal/service"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	RouterConfig *router.Config `yaml:"router" env-prefix:"ROUTER_"`
	DBConfig     *db.Config     `yaml:"db" env-prefix:"DB_"`
	RedisConfig  *redis.Config  `yaml:"redis" env-prefix:"REDIS_"`
	SenderConfig *sender.Config `yaml:"sender" env-prefix:"SENDER_"`
}

func readConfig(filename string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(filename, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	cfg, err := readConfig("./cfg.yml")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Config file read successfully")
	db, err := db.New(cfg.DBConfig)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("DB connected successfully")
	redis, err := redis.New(cfg.RedisConfig)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Redis db connected successfully")
	sender, err := sender.New(cfg.SenderConfig)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Email sender connected successfully")
	service := service.New(redis, db, sender)
	router, _ := router.New(cfg.RouterConfig, service)
	err = router.Listen()
	if err != nil {
		log.Fatalln("Failed to host router:", err.Error())
	}
	log.Printf("Router is listening...")
}
