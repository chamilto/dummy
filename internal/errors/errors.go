package errors

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
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
