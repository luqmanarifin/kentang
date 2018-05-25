package service

import "github.com/go-redis/redis"

type Redis struct {
	db *redis.Client
}

// Option holds all necessary options for Redis.
type Option struct {
	Host     string
	Port     string
	Password string
	Database string
}

func NewRedis(opt Option) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr: opt.Host+opt.Port,
		Password: opt.Password
		DB: opt.Database,
	})
	pong, err := client.Ping().Result()
	log.Println(pong, err)
	return &Redis{db: client}
}
