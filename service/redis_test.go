package service

import (
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

func getRedisConnection(t *testing.T) (*Redis, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	database, _ := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	opt := RedisOption{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Database: database,
	}
	return NewRedis(opt), nil
}

func TestGetKeyword(t *testing.T) {
	r, _ := getRedisConnection(t)
	val, _ := r.GetKeyword("source", "kentang")
	t.Logf("val %d", val)
}

func TestAddKeyword(t *testing.T) {
	r, _ := getRedisConnection(t)
	r.AddKeyword("source", "lala", "1")
}

func TestRemoveAllKeyword(t *testing.T) {
	r, _ := getRedisConnection(t)
	r.RemoveAllKeyword("source")
}
