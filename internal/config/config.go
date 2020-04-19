package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type DBConf struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewDBConf() *DBConf {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		logrus.Fatalf("REDIS_DB must be an integer: %v", err)
	}

	return &DBConf{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	}
}

type Config struct {
	Bind     string
	LogLevel string
	DB       *DBConf
}

func LoadEnv() {
	env := os.Getenv("DUMMY_ENV")

	if env == "" || strings.ToLower(env) == "development" {
		err := godotenv.Load(".env")

		if err != nil {
			logrus.Fatalf("Error loading .env file: %s", err.Error())
		}

	}
}

func NewConfig() *Config {
	return &Config{
		Bind:     os.Getenv("BIND_ADDRESS"),
		LogLevel: os.Getenv("LOG_LEVEL"),
		DB:       NewDBConf(),
	}
}
