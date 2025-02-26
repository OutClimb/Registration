package app

import (
	"github.com/OutClimb/Registration/internal/store"
)

type AppLayer interface {
	AuthenticateUser(username string, password string) (*UserInternal, error)
	CreateToken(user *UserInternal, clientIp string) (string, error)
	CreateSubmission(slug string, ipAddress string, userAgent string, values map[string]string) (*SubmissionInternal, error)
	GetForm(slug string) (*FormInternal, error)
	ValidateToken(userId uint, clientIp string) error
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
