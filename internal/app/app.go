package app

import (
	"github.com/OutClimb/Registration/internal/store"
)

type AppLayer interface {
	AuthenticateUser(username string, password string) (*UserInternal, error)
	CheckRole(userRole string, requiredRole string) bool
	CreateToken(user *UserInternal, clientIp string) (string, error)
	CreateSubmission(slug string, ipAddress string, userAgent string, values map[string]string) (*SubmissionInternal, error)
	GetForm(slug string) (*FormInternal, error)
	GetFormsForUser(userId uint) (*[]FormInternal, error)
	GetSubmissionsForForm(slug string, userId uint) (*[]SubmissionsInternal, error)
	ValidateUser(userId uint) error
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
