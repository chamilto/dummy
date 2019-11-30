package testutil

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/chamilto/dummy/internal/config"
	"github.com/chamilto/dummy/internal/db"
)

func NewRequest(t *testing.T, method, path string) *http.Request {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Fatal(err)
	}

	return req
}

func ValidateStatus(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
	if status := rr.Code; status != expected {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expected)
	}
}

func GetResponseBytes(t *testing.T, rr *httptest.ResponseRecorder) []byte {
	b, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func AssertEquals(t *testing.T, rr *httptest.ResponseRecorder, got, expected interface{}) {
	eq := reflect.DeepEqual(got, expected)

	if !eq {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func NewTestConf(t *testing.T) config.Config {
	c := config.Config{DB: config.DBConf{}}
	c.DB.Host = os.Getenv("TEST_REDIS_HOST")
	c.DB.Port = os.Getenv("TEST_REDIS_PORT")
	c.DB.Password = os.Getenv("TEST_REDIS_PASSWORD")
	db, _ := strconv.Atoi(os.Getenv("TEST_REDIS_DB"))
	c.DB.DB = db

	return c
}
