package app

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/chamilto/dummy/internal/config"
	"github.com/chamilto/dummy/internal/db"
)

type App struct {
	Server *http.Server
	DB     *db.DB
	Router *mux.Router
	Config config.Config
}

func (a *App) Initialize(c config.Config) {
	a.Config = c
	a.DB = db.NewDB(c)
	a.Router = mux.NewRouter()
	a.Server = &http.Server{
		Handler:      a.Router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

}
