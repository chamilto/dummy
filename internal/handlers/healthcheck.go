package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthCheckResponse struct {
	Ok bool `json:"ok"`
}

func (c *HandlerController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(HealthCheckResponse{Ok: true})
}
