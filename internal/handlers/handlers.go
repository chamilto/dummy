package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/chamilto/dummy/internal/config"
	"github.com/chamilto/dummy/internal/db"
)

type HandlerContext struct {
	DB     *db.DB
	Router *mux.Router
	Config *config.Config
}

func NewHandlerContext(c *config.Config, db *db.DB) *HandlerContext {
	return &HandlerContext{
		DB:     db,
		Config: c,
	}
}

type RequestHandlerFunction func(w http.ResponseWriter, r *http.Request)

func HandleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}

func HandleConfigRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler(w, r)
	}
}

func NewRouter(c *HandlerContext) *mux.Router {
	r := mux.NewRouter()
	c.Router = r

	// dummy config api
	r.HandleFunc("/dummy-config/healthcheck", HandleConfigRequest(c.HealthCheck))
	r.HandleFunc(
		"/dummy-config/endpoints",
		HandleConfigRequest(c.CreateDummyEndpoint),
	).Methods("POST")
	r.HandleFunc(
		"/dummy-config/endpoints",
		HandleConfigRequest(c.GetAllDummyEndpoints),
	).Methods("GET")
	r.HandleFunc(
		"/dummy-config/endpoints/{name}",
		HandleConfigRequest(c.GetDetailDummyEndpoint),
	).Methods("GET")
	r.HandleFunc(
		"/dummy-config/endpoints/{name}",
		HandleConfigRequest(c.UpdateDummyEndpoint),
	).Methods("PUT")

	// Hijack the 404 handler to register our Dummy Endpoint matcher
	r.NotFoundHandler = HandleRequest(c.Dummy)

	return r
}
