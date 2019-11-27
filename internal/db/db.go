package db

import (
	"strings"

	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
)

const REDIS_KEY_PREFIX = "dummy"

type DB struct {
	*redis.Client
}

func NewDB() *DB {
	db := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
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
