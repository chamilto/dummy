package handlers

import (
	"encoding/json"
	"github.com/chamilto/dummy/app/dummyendpoint"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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

func WriteServerError(w http.ResponseWriter, msg string, err error) {
	logrus.Warn(msg)
	http.Error(w, err.Error(), 500)

}

func CreateDummyEndpoint(db *redis.Client, w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		WriteServerError(w, "unable to parse request body.", err)
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
		WriteServerError(w, "unable to unmarshal request body to json", err)
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
		WriteServerError(w, "error saving new endpoint to DB", err)
		return
	}

}

func GetAllDummyEndpoints(db *redis.Client, w http.ResponseWriter, r *http.Request) {
	endpoints, err := dummyendpoint.GetAllDummyEndpoints(db)

	if err != nil {
		WriteServerError(w, "unable to fetch dummy endpoints from db", err)
		return
	}

	ret := []dummyendpoint.DummyEndpoint{}

	for _, v := range endpoints {
		e := dummyendpoint.DummyEndpoint{}
		json.Unmarshal([]byte(v), &e)
		ret = append(ret, e)
	}

	json.NewEncoder(w).Encode(ret)
}

func GetDetailDummyEndpoint(db *redis.Client, w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	de, err := dummyendpoint.LoadFromName(db, name)

	if err != nil {
		WriteServerError(w, "unable to fetch dummy endpoint from db", err)
		return

	}

	if de == nil {
		WriteError(
			w,
			"NotFoundError", "Dummy Endpoint not found",
			http.StatusNotFound,
		)
		return
	}
	json.NewEncoder(w).Encode(de)
}

func UpdateDummyEndpoint(db *redis.Client, w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		WriteServerError(w, "unable to parse request body", err)
		return
	}

	valid, validationErrs := dummyendpoint.Validate(b, dummyendpoint.DummyEndpointSchemaLoader)

	if !valid {
		// return validation errors to user
		msg := strings.Join(validationErrs, ",")
		WriteError(w, "ValidationError", msg, http.StatusBadRequest)
		return
	}

	updatedEndpoint := dummyendpoint.DummyEndpoint{}

	var existingEndpoint *dummyendpoint.DummyEndpoint
	existingEndpoint, err = dummyendpoint.LoadFromName(db, mux.Vars(r)["name"])
	WriteServerError(w, "unable to load dummy endpoint", err)

	err = json.Unmarshal(b, &updatedEndpoint)

	if err != nil {
		WriteServerError(w, "unable to unmarshal request body to json", err)

	}

	// make sure the name and pattern exist
	unq, _ := updatedEndpoint.IsUnique(db)

	if unq {
		WriteError(
			w,
			"NotFoundError", "Dummy Endpoint not found",
			http.StatusNotFound,
		)
		return
	}

	// cannot change name
	if existingEndpoint.Name != updatedEndpoint.Name {
		WriteError(
			w,
			"BadRequestError", "Field name cannot be changed.",
			http.StatusNotFound,
		)
		return
	}

	// cannot change path pattern
	if existingEndpoint.PathPattern != updatedEndpoint.PathPattern {
		WriteError(
			w,
			"BadRequestError", "Field pathPattern cannot be changed.",
			http.StatusNotFound,
		)
		return
	}

	err = updatedEndpoint.Save(db)

	if err != nil {
		WriteServerError(w, "unable to save endpoint", err)
		return

	}
}

// Match the incoming request's url path + Method to a dummy endpoint
// Use the dummy endpoint struct data to build our custom response
func Dummy(db *redis.Client, w http.ResponseWriter, r *http.Request) {
	de, err := dummyendpoint.MatchEndpoint(db, r)

	if err != nil {
		WriteServerError(w, "unable to load dummy endpoint", err)
		return

	}

	if de == nil {
		WriteError(
			w,
			"NotFoundError", "URL path + method do not match any existing dummy endpoints.",
			http.StatusNotFound,
		)
		return

	}

	de.SetResponseData(w)
}

func HealthCheck(db *redis.Client, w http.ResponseWriter, r *http.Request) {
	_, err := db.Ping().Result()

	ok := true

	if err != nil {
		ok = false
	}

	json.NewEncoder(w).Encode(map[string]bool{"ok": ok})
}
