package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/chamilto/dummy/internal/db"
)

type RequestHandlerFunction func(db *db.DB, w http.ResponseWriter, r *http.Request)

func HandleRequest(handler RequestHandlerFunction, db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(db, w, r)
	}
}

func HandleConfigRequest(handler RequestHandlerFunction, db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler(db, w, r)
	}
}

func RegisterHandlers(r *mux.Router, db *db.DB) {
	// dummy config api
	r.HandleFunc("/dummy-config/health", HandleConfigRequest(HealthCheck, db))
	r.HandleFunc(
		"/dummy-config/endpoints",
		HandleConfigRequest(CreateDummyEndpoint, db),
	).Methods("POST")
	r.HandleFunc(
		"/dummy-config/endpoints",
		HandleConfigRequest(GetAllDummyEndpoints, db),
	).Methods("GET")
	r.HandleFunc(
		"/dummy-config/endpoints/{name}",
		HandleConfigRequest(GetDetailDummyEndpoint, db),
	).Methods("GET")
	r.HandleFunc(
		"/dummy-config/endpoints/{name}",
		HandleConfigRequest(UpdateDummyEndpoint, db),
	).Methods("PUT")

	// Hijack the 404 handler to register our Dummy Endpoint matcher
	r.NotFoundHandler = HandleRequest(Dummy, db)

}
