package http

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

func (h *httpLayer) getForm(c *gin.Context) {
	form, error := h.app.GetForm(c.Param("slug"))
	if error != nil {
		c.Redirect(http.StatusTemporaryRedirect, "https://outclimb.gay")
		return
	}

	if form.IsBeforeFormOpen() {
		if tmpl, error := template.New("notOpen").ParseFiles("./web/notOpen.html.tmpl"); error != nil {
			c.Status(http.StatusInternalServerError)
			return
		} else if tmpl.Execute(c.Writer, nil); error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
		return
	}

	if form.IsAfterFormClose() {
		if tmpl, error := template.New("closed").ParseFiles("./web/closed.html.tmpl"); error != nil {
			c.Status(http.StatusInternalServerError)
			return
		} else if tmpl.Execute(c.Writer, nil); error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
		return
	}

	if form.IsFormFilled() {
		if tmpl, error := template.New("filled").ParseFiles("./web/filled.html.tmpl"); error != nil {
			c.Status(http.StatusInternalServerError)
			return
		} else if tmpl.Execute(c.Writer, nil); error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
		return
	}

	error = form.WriteTemplate(c.Writer)
	if error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
}
