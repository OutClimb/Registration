package app

import (
	"io"
	"text/template"
)

func (a *appLayer) FormExists(slug string) bool {
	return a.store.GetFormExists(slug)
}

func (a *appLayer) WriteFormTemplate(slug string, writer io.Writer) error {
	if form, error := a.store.GetForm(slug); error != nil {
		return error
	} else if tmpl, error := template.New(form.Template).ParseFiles("./web/" + form.Template + ".html.tmpl"); error != nil {
		return error
	} else if tmpl.Execute(writer, form); error != nil {
		return error
	}

	return nil
}
