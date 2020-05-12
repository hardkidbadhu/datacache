package conn

import (
	"workerservice/config"

	"log"
	"time"

	"github.com/go-redis/redis"
)

// connectRedis connects to redis instance
func ConnectRedis(db int) *redis.Client {
	log.Println("INFO: Connectiong to redis - ", db)
	// Creating base client

	opts := redis.Options{
		Addr:            config.Cfg.RedisVars.ConnString,
		Password:        ``, // no password set
		DB:              db,
		MaxRetries:      config.Cfg.RedisVars.MaxRetries,
		MinRetryBackoff: time.Duration(config.Cfg.RedisVars.MinRetryBackoff) * time.Millisecond,
		DialTimeout:     time.Duration(config.Cfg.RedisVars.DialTimeout) * time.Second,
		ReadTimeout:     time.Duration(config.Cfg.RedisVars.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(config.Cfg.RedisVars.WriteTimeout) * time.Second,
	}

	if config.Cfg.RedisVars.Password != "" {
		opts.Password = config.Cfg.RedisVars.Password
	}

	client := redis.NewClient(&opts)

	return client
}
