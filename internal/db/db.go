package db

import (
	"strings"

	"github.com/chamilto/dummy/internal/config"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
)

const REDIS_KEY_PREFIX = "dummy"

type DB struct {
	*redis.Client
}

func NewDB(c config.Config) *DB {
	db := redis.NewClient(&redis.Options{
		Addr:     c.DB.Host + ":" + c.DB.Port,
		Password: c.DB.Password,
		DB:       c.DB.DB,
	})

	_, err := db.Ping().Result()

	if err != nil {
		logrus.Fatal("unable to connect to redis.")
	}

	return &DB{db}
}

func (_ *DB) BuildKey(parts []string) string {
	// prepend
	parts = append([]string{REDIS_KEY_PREFIX}, parts...)
	return strings.Join(parts, ":")
}
