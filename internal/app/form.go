package app

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/OutClimb/Registration/internal/store"
)

type FormFieldInternal struct {
	Name       string
	Slug       string
	Type       string
	Metadata   *interface{}
	Required   bool
	Validation string
	Order      uint
}

func (f *FormFieldInternal) Internalize(field *store.FormField) {
	f.Name = field.Name
	f.Slug = field.Slug
	f.Type = field.Type
	if field.Metadata.Valid && len(field.Metadata.String) > 0 {
		err := json.Unmarshal([]byte(field.Metadata.String), &f.Metadata)
		if err != nil {
			f.Metadata = nil
		}
	}
	f.Required = field.Required
	if field.Validation.Valid && len(field.Validation.String) > 0 {
		f.Validation = field.Validation.String
	} else {
		f.Validation = ""
	}
	f.Order = field.Order
}

type FormInternal struct {
	Slug             string
	Name             string
	Template         string
	OpensOn          *time.Time
	ClosesOn         *time.Time
	Submissions      uint
	MaxSubmissions   uint
	Fields           map[string]*FormFieldInternal
	RecaptchaSiteKey string
	NotOpenMessage   string
	ClosedMessage    string
	SuccessMessage   string
}

func (f *FormInternal) Internalize(form *store.Form, fields *[]store.FormField, submissions int64) {
	f.Name = form.Name
	f.Slug = form.Slug
	f.Template = form.Template
	f.Submissions = uint(submissions)
	f.MaxSubmissions = form.MaxSubmissions

	f.OpensOn = nil
	if form.OpensOn.Valid {
		f.OpensOn = &form.OpensOn.Time
	}

	f.ClosesOn = nil
	if form.ClosesOn.Valid {
		f.ClosesOn = &form.ClosesOn.Time
	}

	f.NotOpenMessage = "The event is not open for registration just yet, but check back soon!"
	if form.NotOpenMessage.Valid {
		f.NotOpenMessage = form.NotOpenMessage.String
	}

	f.ClosedMessage = "The event is closed for registration. Please check back later for more events."
	if form.ClosedMessage.Valid {
		f.ClosedMessage = form.ClosedMessage.String
	}

	f.SuccessMessage = "Thank you for registering! We'll see you at the event."
	if form.SuccessMessage.Valid {
		f.SuccessMessage = form.SuccessMessage.String
	}

	if fields != nil {
		f.Fields = make(map[string]*FormFieldInternal, len(*fields))
		for _, field := range *fields {
			f.Fields[field.Slug] = &FormFieldInternal{}
			f.Fields[field.Slug].Internalize(&field)
		}
	}
}

func (f *FormInternal) IsBeforeFormOpen() bool {
	if f.OpensOn == nil {
		return false
	}

	return f.OpensOn.After(time.Now())
}

func (f *FormInternal) IsAfterFormClose() bool {
	if f.ClosesOn == nil {
		return false
	}

	return f.ClosesOn.Before(time.Now())
}

func (f *FormInternal) IsFormFilled() bool {
	return f.MaxSubmissions != 0 && f.Submissions >= f.MaxSubmissions
}

func (a *appLayer) GetForm(slug string) (*FormInternal, error) {
	if form, error := a.store.GetForm(slug); error != nil {
		return &FormInternal{}, error
	} else if formFields, error := a.store.GetFormFields(form.ID); error != nil {
		return &FormInternal{}, error
	} else if submissions, error := a.store.GetNumberOfSubmissions(form.ID); error != nil {
		return &FormInternal{}, error
	} else {
		internalForm := FormInternal{}
		internalForm.Internalize(form, formFields, submissions)

		if siteKey, siteKeyExist := os.LookupEnv("RECAPTCHA_SITE_KEY"); siteKeyExist {
			internalForm.RecaptchaSiteKey = siteKey
		}

		return &internalForm, nil
	}
}

func (a *appLayer) GetFormsForUser(userId uint) (*[]FormInternal, error) {
	user, err := a.store.GetUser(userId)
	if err != nil {
		return nil, errors.New("User not found")
	}

	if user.Role == "admin" {
		if forms, err := a.store.GetAllForms(); err != nil {
			return nil, errors.New("Failed to get forms")
		} else {
			internalForms := make([]FormInternal, len(*forms))
			for i, form := range *forms {
				internalForm := FormInternal{}
				internalForm.Internalize(&form, nil, 0)
				internalForms[i] = internalForm
			}

			return &internalForms, nil
		}
	} else if user.Role == "viewer" {
		if forms, err := a.store.GetFormsForUser(userId); err != nil {
			return nil, errors.New("Failed to get forms")
		} else {
			internalForms := make([]FormInternal, len(*forms))
			for i, form := range *forms {
				internalForm := FormInternal{}
				internalForm.Internalize(&form, nil, 0)
				internalForms[i] = internalForm
			}

			return &internalForms, nil
		}
	}

	return nil, errors.New("Unauthorized")
}
