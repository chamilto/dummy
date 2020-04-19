package errors

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	ValidationError = "Validation Error"
	ConflictError   = "Conflict Error"
	NotFoundError   = "Not Found Error"
)

type ErrorMessage struct {
	ErrorType string
	Msg       string
}

func WriteError(w http.ResponseWriter, errType string, msg string, status int) {
	errMsg := ErrorMessage{ErrorType: errType, Msg: msg}
	b, _ := json.Marshal(errMsg)
	w.WriteHeader(status)
	w.Write(b)
}

func WriteServerError(w http.ResponseWriter, msg string, err error) {
	logrus.Error(msg)
	http.Error(w, err.Error(), 500)
}
