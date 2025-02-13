package app

import (
	"time"

	"github.com/OutClimb/Registration/internal/store"
)

type FormFieldInternal struct {
	Name       string
	Slug       string
	Type       string
	Metadata   string
	Required   bool
	Validation string
}

func (f *FormFieldInternal) Internalize(field store.FormField) {
	f.Name = field.Name
	f.Slug = field.Slug
	f.Type = field.Type
	f.Metadata = field.Metadata
	f.Required = field.Required
	f.Validation = field.Validation
}

type FormInternal struct {
	Slug           string
	Name           string
	Template       string
	OpensOn        *time.Time
	ClosesOn       *time.Time
	Submissions    uint
	MaxSubmissions uint
	Fields         map[string]*FormFieldInternal
}

func (f *FormInternal) Internalize(field store.Form) {
	f.Name = field.Name
	f.Slug = field.Slug
	f.Template = field.Template
	f.OpensOn = field.OpensOn
	f.ClosesOn = field.ClosesOn
	f.Submissions = uint(len(field.Submissions))
	f.MaxSubmissions = field.MaxSubmissions

	f.Fields = make(map[string]*FormFieldInternal, len(field.FormFields))
	for _, field := range field.FormFields {
		f.Fields[field.Slug] = &FormFieldInternal{}
		f.Fields[field.Slug].Internalize(field)
	}
}

func (f *FormInternal) IsBeforeFormOpen() bool {
	if f.OpensOn == nil {
		return false
	}

	return f.OpensOn.Before(time.Now())
}

func (f *FormInternal) IsAfterFormClose() bool {
	if f.ClosesOn == nil {
		return false
	}

	return f.ClosesOn.After(time.Now())
}

func (f *FormInternal) IsFormFilled() bool {
	return f.MaxSubmissions != 0 && f.Submissions >= f.MaxSubmissions
}

func (a *appLayer) GetForm(slug string) (FormInternal, error) {
	form, error := a.store.GetForm(slug)
	if error != nil {
		return FormInternal{}, error
	}

	internalForm := FormInternal{}
	internalForm.Internalize(form)

	return internalForm, nil
}
