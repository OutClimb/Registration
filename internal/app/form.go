package app

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
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
	if field.Metadata != nil && len(*field.Metadata) > 0 {
		err := json.Unmarshal([]byte(*field.Metadata), &f.Metadata)
		if err != nil {
			f.Metadata = nil
		}
	}
	f.Required = field.Required
	if field.Validation != nil {
		f.Validation = *field.Validation
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
	f.OpensOn = form.OpensOn
	f.ClosesOn = form.ClosesOn
	f.Submissions = uint(submissions)
	f.MaxSubmissions = form.MaxSubmissions

	f.NotOpenMessage = "The event is not open for registration just yet, but check back soon!"
	if form.NotOpenMessage != nil {
		f.NotOpenMessage = *form.NotOpenMessage
	}

	f.ClosedMessage = "The event is closed for registration. Please check back later for more events."
	if form.ClosedMessage != nil {
		f.ClosedMessage = *form.ClosedMessage
	}

	f.SuccessMessage = "Thank you for registering! We'll see you at the event."
	if form.SuccessMessage != nil {
		f.SuccessMessage = *form.SuccessMessage
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

func (a *appLayer) DuplicateForm(slug string) (*FormInternal, error) {
	originalForm, err := a.store.GetForm(slug)
	if err != nil {
		return nil, err
	}

	originalFormFields, err := a.store.GetFormFields(originalForm.ID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	newForm, err := a.store.CreateForm(
		"Copy of "+originalForm.Name,
		originalForm.Slug+"-"+strconv.FormatInt(now.UnixMilli(), 10),
		originalForm.Template,
		originalForm.OpensOn,
		originalForm.ClosesOn,
		originalForm.MaxSubmissions,
		originalForm.NotOpenMessage,
		originalForm.ClosedMessage,
		originalForm.SuccessMessage,
		originalForm.EmailFormFieldSlug,
		originalForm.EmailTo,
		originalForm.EmailSubject,
		originalForm.EmailTemplate,
	)
	if err != nil {
		return nil, err
	}

	formFields := make([]store.FormField, len(*originalFormFields))
	for i, originalFormField := range *originalFormFields {
		formField, err := a.store.CreateFormField(
			newForm.ID,
			originalFormField.Name,
			originalFormField.Slug,
			originalFormField.Type,
			originalFormField.Metadata,
			originalFormField.Required,
			originalFormField.Validation,
			originalFormField.Order,
		)
		if err != nil {
			// TODO: Clean up form
			return nil, err
		}

		formFields[i] = *formField
	}

	internalForm := FormInternal{}
	internalForm.Internalize(newForm, &formFields, 0)

	return &internalForm, nil
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
