package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	ctlr := newTestingHandlerController()
	rr := httptest.NewRecorder()
	req := newRequest(t, "GET", "/dummy-config/healthcheck", nil)
	ctlr.Router.ServeHTTP(rr, req)

	assertStatus(t, rr, http.StatusOK)

	expected := HealthCheckResponse{Ok: true}
	got := HealthCheckResponse{}
	json.Unmarshal(getResponseBytes(t, rr), &got)

	assertEquals(t, got, expected)
}
