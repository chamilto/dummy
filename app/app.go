package app

import (
	"github.com/chamilto/dummy/app/handlers"
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

func (a *App) registerHandlers() {
	// dummy config api
	a.Router.HandleFunc("/dummy-config/health", a.handleRequest(handlers.HealthCheck))
	a.Router.HandleFunc(
		"/dummy-config/endpoints",
		a.handleRequest(handlers.CreateDummyEndpoint),
	).Methods("POST")
	a.Router.HandleFunc(
		"/dummy-config/endpoints",
		a.handleRequest(handlers.GetAllDummyEndpoints),
	).Methods("GET")
	a.Router.HandleFunc(
		"/dummy-config/endpoints/{name}",
		a.handleRequest(handlers.GetDetailDummyEndpoint),
	).Methods("GET")

	// Hijack the 404 handler to register our Dummy Endpoint matcher
	a.Router.NotFoundHandler = a.handleRequest(handlers.Dummy)

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

type RequestHandlerFunction func(db *redis.Client, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}
