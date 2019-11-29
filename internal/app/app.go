package app

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/chamilto/dummy/internal/db"
)

type App struct {
	Server *http.Server
	DB     *db.DB
	Router *mux.Router
}

// todo: take config struct param
func (a *App) Initialize() {
	a.DB = db.NewDB()
	a.Router = mux.NewRouter()
	a.Server = &http.Server{
		Handler:      a.Router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

}
