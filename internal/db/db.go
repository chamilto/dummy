package db

import (
	"fmt"

	"github.com/chamilto/dummy/internal/config"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
)

const REDIS_KEY_PREFIX = "dummy"

type RedisClient interface {
	redis.UniversalClient
	BuildKey(key string) string
}

type DB struct {
	redis.UniversalClient
}

func NewDB(c *config.Config) *DB {
	db := redis.NewClient(&redis.Options{
		Addr:     c.DB.Host + ":" + c.DB.Port,
		Password: c.DB.Password,
		DB:       c.DB.DB,
	})

	_, err := db.Ping().Result()

	if err != nil {
		logrus.Fatal("unable to connect to redis")
	}

	return &DB{db}
}

func (DB) BuildKey(key string) string {
	return fmt.Sprintf("%s:%s", REDIS_KEY_PREFIX, key)
}
