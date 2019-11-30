// +build integration

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chamilto/dummy/internal/app"
	"github.com/chamilto/dummy/internal/testutil"
)

func TestHealthCheckHandler(t *testing.T) {
	a := app.App{}
	c := testutil.NewTestConf(t)
	a.Initialize(c)
	RegisterHandlers(a.Router, a.DB)
	rr := httptest.NewRecorder()
	req := testutil.NewRequest(t, "GET", "/dummy-config/health")
	a.Router.ServeHTTP(rr, req)
	testutil.ValidateStatus(t, rr, http.StatusOK)
	expected := HealthCheckResponse{Ok: true}
	got := HealthCheckResponse{}
	json.Unmarshal(testutil.GetResponseBytes(t, rr), &got)
	testutil.AssertEquals(t, rr, got, expected)
}
