package main

import (
	"github.com/OutClimb/Registration/internal/app"
	"github.com/OutClimb/Registration/internal/http"
	"github.com/OutClimb/Registration/internal/store"
)

func main() {
	storeLayer := store.New()
	appLayer := app.New(storeLayer)
	httpLayer := http.New(appLayer)

	httpLayer.Run()
}
