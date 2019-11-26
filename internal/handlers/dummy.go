package handlers

import (
	"github.com/chamilto/dummy/internal/dummyendpoint"
	"github.com/chamilto/dummy/internal/errors"
	"github.com/go-redis/redis/v7"
	"net/http"
)

// Match the incoming request's url path + Method to a dummy endpoint
// Use the dummy endpoint struct data to build our custom response
func Dummy(db *redis.Client, w http.ResponseWriter, r *http.Request) {
	de, err := dummyendpoint.MatchEndpoint(db, r)

	if err != nil {
		errors.WriteServerError(w, "unable to load dummy endpoint", err)
		return

	}

	if de == nil {
		errors.WriteError(
			w,
			"NotFoundError", "URL path + method do not match any existing dummy endpoints.",
			http.StatusNotFound,
		)
		return

	}

	de.SetResponseData(w)
}
