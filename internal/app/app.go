package app

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/chamilto/dummy/internal/db"
	"github.com/chamilto/dummy/internal/handlers"
)

type App struct {
	Server *http.Server
	DB     *db.DB
	Router *mux.Router
}

func (a *App) registerHandlers() {
	// dummy config api
	a.Router.HandleFunc("/dummy-config/health", a.handleRequest(handlers.HealthCheck))
	a.Router.HandleFunc(
		"/dummy-config/endpoints",
		a.handleConfigRequest(handlers.CreateDummyEndpoint),
	).Methods("POST")
	a.Router.HandleFunc(
		"/dummy-config/endpoints",
		a.handleConfigRequest(handlers.GetAllDummyEndpoints),
	).Methods("GET")
	a.Router.HandleFunc(
		"/dummy-config/endpoints/{name}",
		a.handleConfigRequest(handlers.GetDetailDummyEndpoint),
	).Methods("GET")
	a.Router.HandleFunc(
		"/dummy-config/endpoints/{name}",
		a.handleConfigRequest(handlers.UpdateDummyEndpoint),
	).Methods("PUT")

	// Hijack the 404 handler to register our Dummy Endpoint matcher
	a.Router.NotFoundHandler = a.handleRequest(handlers.Dummy)

}

// todo: take config struct param
func (a *App) Initialize() {
	a.DB = db.NewDB()
	a.Router = mux.NewRouter()
	a.registerHandlers()
	a.Server = &http.Server{
		Handler:      a.Router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

}

type RequestHandlerFunction func(db *db.DB, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}

func (a *App) handleConfigRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler(a.DB, w, r)
	}
}
