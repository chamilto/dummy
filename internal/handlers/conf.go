// Dummy endpoint configuration handlers
package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/chamilto/dummy/internal/errors"
	"github.com/chamilto/dummy/internal/models/dummy"
	"github.com/chamilto/dummy/internal/utils"
)

type collectionResponse struct {
	Items []dummy.DummyEndpoint `json:"items"`
}

func (c *HandlerController) CreateDummyEndpoint(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		errors.WriteServerError(w, "unable to parse request body", err)
		return
	}

	valid, validationErrs := utils.ValidateJson(b, dummy.DummyEndpointSchemaLoader)

	if !valid {
		msg := strings.Join(validationErrs, ",")
		errors.WriteError(w, errors.ValidationError, msg, http.StatusBadRequest)
		return
	}

	var newEndpoint dummy.DummyEndpoint

	err = json.Unmarshal(b, &newEndpoint)

	if err != nil {
		errors.WriteServerError(w, "unable to unmarshal request body to json", err)
		return
	}

	// check name + pattern uniqueness before saving
	unq, unqErrMsg := newEndpoint.IsUnique(c.DB)

	if !unq {
		errors.WriteError(w, errors.ConflictError, unqErrMsg, http.StatusConflict)
		return
	}

	err = newEndpoint.Save(c.DB)

	fmt.Printf("%v", err)

	if err != nil {
		errors.WriteServerError(w, "error saving new endpoint to DB", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *HandlerController) GetAllDummyEndpoints(w http.ResponseWriter, r *http.Request) {
	endpoints, err := dummy.GetAllDummyEndpoints(c.DB)

	if err != nil {
		errors.WriteServerError(w, "unable to fetch dummy endpoints from db", err)
		return
	}

	items := []dummy.DummyEndpoint{}

	for _, v := range endpoints {
		de := dummy.DummyEndpoint{}
		json.Unmarshal([]byte(v), &de)
		items = append(items, de)
	}

	ret := collectionResponse{Items: items}

	b, err := json.Marshal(&ret)

	if err != nil {
		errors.WriteServerError(w, "error encoding dummy endpoints", err)
		return
	}

	w.Write(b)
}

func (c *HandlerController) GetDetailDummyEndpoint(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	de, err := dummy.LoadFromName(c.DB, name)

	if err != nil {
		errors.WriteServerError(w, "unable to fetch dummy endpoint from db", err)
		return

	}

	if de == nil {
		errors.WriteError(
			w,
			errors.NotFoundError, "Dummy Endpoint not found",
			http.StatusNotFound,
		)
		return
	}

	b, err := de.ToJSON() // maybe dumb?

	if err != nil {
		errors.WriteServerError(w, "unable to encode dummy endpoint", err)
		return
	}

	w.Write(b)
}

func (c *HandlerController) UpdateDummyEndpoint(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		errors.WriteServerError(w, "unable to parse request body", err)
		return
	}

	valid, validationErrs := utils.ValidateJson(b, dummy.DummyEndpointSchemaLoader)

	if !valid {
		msg := strings.Join(validationErrs, ",")
		errors.WriteError(w, errors.ValidationError, msg, http.StatusBadRequest)
		return
	}

	updatedEndpoint := dummy.DummyEndpoint{}

	existingEndpoint, err := dummy.LoadFromName(c.DB, mux.Vars(r)["name"])

	if err != nil {
		errors.WriteServerError(w, "unable to load dummy endpoint", err)
		return
	}

	err = json.Unmarshal(b, &updatedEndpoint)

	if err != nil {
		errors.WriteServerError(w, "unable to unmarshal request body to json", err)

	}

	// make sure the name and pattern exist
	unq, _ := updatedEndpoint.IsUnique(c.DB)

	if unq {
		errors.WriteError(
			w,
			errors.NotFoundError, "Dummy Endpoint not found",
			http.StatusNotFound,
		)
		return
	}

	// cannot change name
	if existingEndpoint.Name != updatedEndpoint.Name {
		errors.WriteError(
			w,
			errors.ValidationError, "Field name cannot be changed.",
			http.StatusNotFound,
		)
		return
	}

	// cannot change path pattern
	if existingEndpoint.PathPattern != updatedEndpoint.PathPattern {
		errors.WriteError(
			w,
			errors.ValidationError, "Field pathPattern cannot be changed.",
			http.StatusNotFound,
		)
		return
	}

	err = updatedEndpoint.Save(c.DB)

	if err != nil {
		errors.WriteServerError(w, "unable to save endpoint", err)
		return

	}
}
