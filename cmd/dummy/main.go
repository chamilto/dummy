package main

import (
	"github.com/chamilto/dummy/internal/app"
)

func main() {
	a := app.App{}
	a.Initialize()
	a.Server.ListenAndServe()
}
