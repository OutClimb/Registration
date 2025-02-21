package app

import (
	"github.com/OutClimb/Registration/internal/store"
)

type AppLayer interface {
	CreateSubmission(slug string, ipAddress string, userAgent string, values map[string]string) (*SubmissionInternal, error)
	GetForm(slug string) (*FormInternal, error)
	ValidateRecaptchaToken(token string, clientIp string) error
	ValidateSubmissionWithForm(submission map[string]string, form *FormInternal) []error
}

type appLayer struct {
	store store.StoreLayer
}

func New(storeLayer store.StoreLayer) *appLayer {
	return &appLayer{
		store: storeLayer,
	}
}
