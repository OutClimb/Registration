package app

import (
	"io"
	"text/template"
	"time"
)

type FormInternal struct {
	ID             uint
	Slug           string
	Name           string
	Template       string
	OpensOn        time.Time
	ClosesOn       time.Time
	Submissions    uint
	MaxSubmissions uint
}

func (f *FormInternal) IsBeforeFormOpen() bool {
	return f.OpensOn.Before(time.Now())
}

func (f *FormInternal) IsAfterFormClose() bool {
	return f.ClosesOn.After(time.Now())
}

func (f *FormInternal) IsFormFilled() bool {
	return f.Submissions >= f.MaxSubmissions
}

func (f *FormInternal) WriteTemplate(writer io.Writer) error {
	if tmpl, error := template.New(f.Template).ParseFiles("./web/" + f.Template + ".html.tmpl"); error != nil {
		return error
	} else if tmpl.Execute(writer, f); error != nil {
		return error
	}

	return nil
}

func (a *appLayer) GetForm(slug string) (FormInternal, error) {
	form, error := a.store.GetForm(slug)
	if error != nil {
		return FormInternal{}, error
	}

	return FormInternal{
		ID:             form.ID,
		Slug:           form.Slug,
		Name:           form.Name,
		Template:       form.Template,
		OpensOn:        form.OpensOn,
		ClosesOn:       form.ClosesOn,
		Submissions:    uint(len(form.Submissions)),
		MaxSubmissions: form.MaxSubmissions,
	}, nil
}
