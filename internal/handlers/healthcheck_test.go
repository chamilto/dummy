package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chamilto/dummy/internal/app"
)

func TestHealthCheckHandler(t *testing.T) {
	a := app.App{}
	a.Initialize()
	RegisterHandlers(a.Router, a.DB)
	rr := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/dummy-config/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	a.Router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := HealthCheckResponse{Ok: true}
	got := HealthCheckResponse{}
	var b []byte

	b, err = ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	json.Unmarshal(b, &got)

	if got != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
