package app

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type App struct {
	Server *http.Server
	DB     *redis.Client
	Router *mux.Router
}

// register all dummy config api handlers
func (a *App) registerHandlers() {
	a.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

}

func getRedisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisClient.Ping().Result()

	if err != nil {
		panic(err)
	}

	return redisClient
}

// todo: take config struct param
func (a *App) Initialize() {
	a.DB = getRedisClient()
	a.Router = mux.NewRouter()
	a.registerHandlers()
	a.Server = &http.Server{
		Handler:      a.Router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

}
