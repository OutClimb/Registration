package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
