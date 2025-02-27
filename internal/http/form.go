package http

import (
	"net/http"
	"time"

	"github.com/OutClimb/Registration/internal/app"
	"github.com/gin-gonic/gin"
)

type FormPublic struct {
	Slug           string `json:"slug"`
	Name           string `json:"name"`
	Template       string `json:"template"`
	OpensOn        string `json:"opens_on"`
	ClosesOn       string `json:"closes_on"`
	MaxSubmissions uint   `json:"max_submissions"`
}

func (f *FormPublic) Publicize(form *app.FormInternal) {
	f.Slug = form.Slug
	f.Name = form.Name
	f.Template = form.Template
	f.MaxSubmissions = form.MaxSubmissions

	if form.OpensOn != nil {
		f.OpensOn = form.OpensOn.Format(time.UnixDate)
	}

	if form.ClosesOn != nil {
		f.ClosesOn = form.ClosesOn.Format(time.UnixDate)
	}
}

func (h *httpLayer) getForm(c *gin.Context) {
	form, error := h.app.GetForm(c.Param("slug"))
	if error != nil {
		c.Redirect(http.StatusTemporaryRedirect, "https://outclimb.gay")
		return
	}

	if form.IsBeforeFormOpen() {
		c.HTML(http.StatusOK, "notOpen.html.tmpl", form)
		return
	}

	if form.IsAfterFormClose() {
		c.HTML(http.StatusOK, "closed.html.tmpl", form)
		return
	}

	if form.IsFormFilled() {
		c.HTML(http.StatusOK, "filled.html.tmpl", form)
		return
	}

	c.HTML(http.StatusOK, form.Template+".html.tmpl", form)
}

func (h *httpLayer) getFormsApi(c *gin.Context) {
	userId := c.GetUint("user_id")

	if internalForms, err := h.app.GetFormsForUser(userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	} else {
		forms := make([]FormPublic, len(*internalForms))
		for i, form := range *internalForms {
			forms[i].Publicize(&form)
		}

		c.JSON(http.StatusOK, forms)
	}
}
