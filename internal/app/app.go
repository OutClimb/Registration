package app

import (
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/OutClimb/Registration/internal/store"
)

type AppLayer interface {
	AuthenticateUser(username string, password string) (*UserInternal, error)
	CheckRole(userRole string, requiredRole string) bool
	CreateSubmission(slug string, ipAddress string, userAgent string, values map[string]string) (*SubmissionInternal, error)
	DeleteSubmission(id uint) error
	DuplicateForm(slug string) (*FormInternal, error)
	GetForm(slug string) (*FormInternal, error)
	GetFormsForUser(userId uint) (*[]FormInternal, error)
	GetSubmissionsForForm(slug string, userId uint) (*[]SubmissionsInternal, error)
	GetUser(userId uint) (*UserInternal, error)
	ValidatePassword(user *UserInternal, password string) error
	ValidateUser(userId uint) error
	ValidateRecaptchaToken(token string, clientIp string) error
	ValidateSubmissionWithForm(submission map[string]string, form *FormInternal) []error
	UpdatePassword(user *UserInternal, password string) error
}

type appLayer struct {
	store          store.StoreLayer
	emailTemplates map[string]*template.Template
}

func New(storeLayer store.StoreLayer) *appLayer {
	matches, _ := filepath.Glob("web/emails/*.html.tmpl")
	emailTemplates := make(map[string]*template.Template, len(matches))

	for _, match := range matches {
		fmt.Println(match)
		name := match[11:(len(match) - 10)]
		template, err := template.ParseFiles(match)

		if err == nil {
			emailTemplates[name] = template
		}
	}

	return &appLayer{
		store:          storeLayer,
		emailTemplates: emailTemplates,
	}
}
