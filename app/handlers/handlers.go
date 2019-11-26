package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/chamilto/dummy/app/dummyendpoint"
	"github.com/go-redis/redis/v7"
	"io/ioutil"
	"net/http"
	"strings"
)

type ErrorMessage struct {
	ErrorType string
	Msg       string
}

func WriteError(w http.ResponseWriter, errType string, msg string, status int) {
	errMsg := ErrorMessage{ErrorType: errType, Msg: msg}
	errB, _ := json.Marshal(errMsg)
	w.WriteHeader(status)
	w.Write(errB)
}

func CreateDummyEndpoint(db *redis.Client, w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		// TODO: logging
		http.Error(w, err.Error(), 500)
		return
	}

	valid, validationErrs := dummyendpoint.Validate(b, dummyendpoint.DummyEndpointSchemaLoader)

	if !valid {
		// return validation errors to user
		msg := strings.Join(validationErrs, ",")
		WriteError(w, "ValidationError", msg, http.StatusBadRequest)
		return
	}

	newEndpoint := dummyendpoint.DummyEndpoint{}

	err = json.Unmarshal(b, &newEndpoint)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// check name + pattern uniqueness before saving
	unq, unqErrMsg := newEndpoint.IsUnique(db)

	if !unq {
		WriteError(w, "ConflictError", unqErrMsg, http.StatusConflict)
		return
	}

	saveErr := newEndpoint.Save(db)

	if saveErr != nil {
		fmt.Println("Error saving new endpoint to DB.")
		http.Error(w, err.Error(), 500)
		return
	}

}

// Match the incoming request's url path + Method to a dummy endpoint
// Use the dummy endpoint struct data to build our custom response
func Dummy(db *redis.Client, w http.ResponseWriter, r *http.Request) {
	de := dummyendpoint.MatchEndpoint(db, r)

	if de == nil {
		WriteError(
			w,
			"NotFoundError", "URL path + method do not match any existing dummy endpoints.",
			http.StatusNotFound,
		)
		return

	}

	json.NewEncoder(w).Encode(map[string]string{"testing body": de.Body})

}

func HealthCheck(db *redis.Client, w http.ResponseWriter, r *http.Request) {
	_, err := db.Ping().Result()

	ok := true

	if err != nil {
		ok = false
	}

	json.NewEncoder(w).Encode(map[string]bool{"ok": ok})
}
