package db

import (
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
)

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
