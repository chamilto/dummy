package handlers

import (
	"net/http"

	"github.com/chamilto/dummy/internal/errors"
	"github.com/chamilto/dummy/internal/models/dummy"
)

// Match the incoming request's url path + Method to a dummy endpoint
// Use the dummy endpoint struct data to build our custom response
func (c *HandlerController) Dummy(w http.ResponseWriter, r *http.Request) {
	de, err := dummy.MatchEndpoint(c.DB, r)

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
	de.RunDelay()
}
