package app

import (
	"fmt"
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

func (f *FormFieldInternal) Internalize(field *store.FormField) {
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

func (f *FormInternal) Internalize(form *store.Form, fields *[]store.FormField, submissions int64) {
	f.Name = form.Name
	f.Slug = form.Slug
	f.Template = form.Template
	f.OpensOn = form.OpensOn
	f.ClosesOn = form.ClosesOn
	f.Submissions = uint(submissions)
	f.MaxSubmissions = form.MaxSubmissions

	f.Fields = make(map[string]*FormFieldInternal, len(*fields))
	for _, field := range *fields {
		f.Fields[field.Slug] = &FormFieldInternal{}
		f.Fields[field.Slug].Internalize(&field)
	}
}

func (f *FormInternal) IsBeforeFormOpen() bool {
	fmt.Printf("%v\n", f.OpensOn)

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

		return &internalForm, nil
	}
}
