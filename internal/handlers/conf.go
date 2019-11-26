package handlers

import (
	"encoding/json"
	"github.com/chamilto/dummy/internal/dummyendpoint"
	"github.com/chamilto/dummy/internal/errors"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
)

func CreateDummyEndpoint(db *redis.Client, w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		errors.WriteServerError(w, "unable to parse request body.", err)
		return
	}

	valid, validationErrs := dummyendpoint.Validate(b, dummyendpoint.DummyEndpointSchemaLoader)

	if !valid {
		// return validation errors to user
		msg := strings.Join(validationErrs, ",")
		errors.WriteError(w, "ValidationError", msg, http.StatusBadRequest)
		return
	}

	newEndpoint := dummyendpoint.DummyEndpoint{}

	err = json.Unmarshal(b, &newEndpoint)

	if err != nil {
		errors.WriteServerError(w, "unable to unmarshal request body to json", err)
		return
	}

	// check name + pattern uniqueness before saving
	unq, unqErrMsg := newEndpoint.IsUnique(db)

	if !unq {
		errors.WriteError(w, "ConflictError", unqErrMsg, http.StatusConflict)
		return
	}

	saveErr := newEndpoint.Save(db)

	if saveErr != nil {
		errors.WriteServerError(w, "error saving new endpoint to DB", err)
		return
	}

}

func GetAllDummyEndpoints(db *redis.Client, w http.ResponseWriter, r *http.Request) {
	endpoints, err := dummyendpoint.GetAllDummyEndpoints(db)

	if err != nil {
		errors.WriteServerError(w, "unable to fetch dummy endpoints from db", err)
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
		errors.WriteServerError(w, "unable to fetch dummy endpoint from db", err)
		return

	}

	if de == nil {
		errors.WriteError(
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
		errors.WriteServerError(w, "unable to parse request body", err)
		return
	}

	valid, validationErrs := dummyendpoint.Validate(b, dummyendpoint.DummyEndpointSchemaLoader)

	if !valid {
		// return validation errors to user
		msg := strings.Join(validationErrs, ",")
		errors.WriteError(w, "ValidationError", msg, http.StatusBadRequest)
		return
	}

	updatedEndpoint := dummyendpoint.DummyEndpoint{}

	var existingEndpoint *dummyendpoint.DummyEndpoint
	existingEndpoint, err = dummyendpoint.LoadFromName(db, mux.Vars(r)["name"])

	if err != nil {

		errors.WriteServerError(w, "unable to load dummy endpoint", err)
		return
	}

	err = json.Unmarshal(b, &updatedEndpoint)

	if err != nil {
		errors.WriteServerError(w, "unable to unmarshal request body to json", err)

	}

	// make sure the name and pattern exist
	unq, _ := updatedEndpoint.IsUnique(db)

	if unq {
		errors.WriteError(
			w,
			"NotFoundError", "Dummy Endpoint not found",
			http.StatusNotFound,
		)
		return
	}

	// cannot change name
	if existingEndpoint.Name != updatedEndpoint.Name {
		errors.WriteError(
			w,
			"BadRequestError", "Field name cannot be changed.",
			http.StatusNotFound,
		)
		return
	}

	// cannot change path pattern
	if existingEndpoint.PathPattern != updatedEndpoint.PathPattern {
		errors.WriteError(
			w,
			"BadRequestError", "Field pathPattern cannot be changed.",
			http.StatusNotFound,
		)
		return
	}

	err = updatedEndpoint.Save(db)

	if err != nil {
		errors.WriteServerError(w, "unable to save endpoint", err)
		return

	}
}
