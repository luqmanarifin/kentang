package service

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/luqmanarifin/kentang/util"
)

type Redis struct {
	db *redis.Client
}

// Option holds all necessary options for Redis
type RedisOption struct {
	Host     string
	Port     string
	Password string
	Database int
}

func NewRedis(opt RedisOption) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     opt.Host + ":" + opt.Port,
		Password: opt.Password,
		DB:       opt.Database,
	})
	pong, err := client.Ping().Result()
	log.Println(pong, err)
	if err != nil {
		return &Redis{}, err
	}

	log.Printf("Success connecting Redis to %s:%s with pass %s\n", opt.Host, opt.Port, opt.Password)
	return &Redis{db: client}, nil
}

// 0 not exist, 1 exist, -1 don't know
func (r *Redis) GetKeyword(source, keyword string) (string, error) {
	val, err := r.db.Get(source + ":" + keyword).Result()
	if err == redis.Nil {
		return "", err
	}
	return val, nil
}

func (r *Redis) AddKeyword(source, keyword, val string) error {
	return r.db.Set(source+":"+keyword, val, 0).Err()
}

func (r *Redis) RemoveKeyword(source, keyword string) error {
	return r.AddKeyword(source, keyword, util.NOT_EXIST)
}

func (r *Redis) RemoveAllKeyword(source string) error {
	script := "return redis.call('DEL', unpack(redis.call('KEYS', ARGV[1] .. '*')))"
	s := make([]string, 0)
	return r.db.Eval(script, s, source+":").Err()
}

func (r *Redis) GetDisplayName(userId string) (string, error) {
	name, err := r.db.Get(userId).Result()
	if err != nil {
		return "", err
	}
	return name, nil
}

func (r *Redis) SetDisplayName(userId, name string) error {
	return r.db.Set(userId, name, 10*24*time.Hour).Err()
}
