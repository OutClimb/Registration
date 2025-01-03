package app

import "github.com/OutClimb/Registration/internal/store"

type AppLayer interface {
	CheckEventExists(slug string) bool
}

type appLayer struct {
	store store.StoreLayer
}

func New(storeLayer store.StoreLayer) *appLayer {
	return &appLayer{
		store: storeLayer,
	}
}
