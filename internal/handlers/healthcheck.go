package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chamilto/dummy/internal/db"
)

type HealthCheckResponse struct {
	Ok bool `json:"ok"`
}

func HealthCheck(db *db.DB, w http.ResponseWriter, r *http.Request) {
	_, err := db.Ping().Result()

	ok := true

	if err != nil {
		ok = false
	}
	json.NewEncoder(w).Encode(HealthCheckResponse{Ok: ok})
}
