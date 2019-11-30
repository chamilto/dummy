package main

import (
	"github.com/chamilto/dummy/internal/app"
	"github.com/chamilto/dummy/internal/config"
	"github.com/chamilto/dummy/internal/handlers"
)

func main() {
	config.LoadEnv()
	c := config.NewConfig()
	a := app.App{}
	a.Initialize(c)
	handlers.RegisterHandlers(a.Router, a.DB)
	a.Server.ListenAndServe()
}
