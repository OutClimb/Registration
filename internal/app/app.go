package app

import (
	"io"

	"github.com/OutClimb/Registration/internal/store"
)

type AppLayer interface {
	FormExists(slug string) bool
	WriteFormTemplate(slug string, writer io.Writer) error
}

type appLayer struct {
	store store.StoreLayer
}

func New(storeLayer store.StoreLayer) *appLayer {
	return &appLayer{
		store: storeLayer,
	}
}
