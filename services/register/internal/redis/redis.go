package redis

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Host     string `yaml:"host" env-prefix:"HOST"`
	Port     string `yaml:"port" env-prefix:"PORT"`
	Password string `yaml:"password" env-prefix:"PASSWORD"`
	DB       int    `yaml:"dbnum" env-prefix:"DB"`
	TTL      int    `yaml:"ttl"`
}

type Redis struct {
	cfg    *Config
	client *redis.Client
	ttl    time.Duration
}

func New(cfg *Config) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	ret := Redis{client: rdb, cfg: cfg, ttl: time.Duration(cfg.TTL) * time.Second}
	return &ret, nil
}

func (r *Redis) GetCode(email string) (string, error) {
	code, err := r.client.Get(context.Background(), "code_email:"+email).Result()
	if err == redis.Nil {
		slog.Debug(fmt.Sprintf("User %v tried to access expired code", email))
		return "", nil
	}
	if err != nil {
		slog.Error(fmt.Sprintf("redis GetCode error: %v", err))
		return "", err
	}
	return code, nil
}

func (r *Redis) PutCode(email, code string) error {
	err := r.client.Set(context.Background(), "code_email:"+email, code, r.ttl)
	if err.Err() != nil {
		slog.Error(fmt.Sprintf("redis PutCode error: %v", err))
		return err.Err()
	}
	return nil
}

func (r *Redis) DeleteCode(email string) error {
	err := r.client.Del(context.Background(), "code_email:"+email)
	if err.Err() != nil {
		log.Println("Cant delete verification code:", err)
		return err.Err()
	}
	return nil
}
