package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConf struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type Config struct {
	DB DBConf
}

func LoadEnv() {
	env := os.Getenv("DUMMY_ENV")

	if env != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error loading .env file")
		}

	}
}

func NewConfig() Config {
	c := Config{DB: DBConf{}}
	c.DB.Host = os.Getenv("REDIS_HOST")
	c.DB.Port = os.Getenv("REDIS_PORT")
	c.DB.Password = os.Getenv("REDIS_PASSWORD")
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	c.DB.DB = db

	return c
}
