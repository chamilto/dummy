package main

import (
	"github.com/chamilto/dummy/internal/app"
	"github.com/chamilto/dummy/internal/handlers"
)

func main() {
	a := app.App{}
	a.Initialize()
	handlers.RegisterHandlers(a.Router, a.DB)
	a.Server.ListenAndServe()
}
