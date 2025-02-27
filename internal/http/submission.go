package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *httpLayer) createSubmission(c *gin.Context) {
	// Get the form we are submitting to
	form, error := h.app.GetForm(c.Param("slug"))
	if error != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Make sure the form is open
	if form.IsBeforeFormOpen() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The event is not open for registration just yet, but check back soon!"})
		return
	}

	// Make sure the form is not closed
	if form.IsAfterFormClose() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The event is closed for registration. Please check back later for more events."})
		return
	}

	// Make sure the form is not full
	if form.IsFormFilled() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The event is full. Please check back later for more events."})
		return
	}

	// Get the submission data
	bodyAsByteArray, err := c.GetRawData()
	if err != nil {
		fmt.Printf("Error: Retrieving request body (%s)\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Parse the submission data
	jsonMap := make(map[string]string)
	err = json.Unmarshal(bodyAsByteArray, &jsonMap)
	if err != nil {
		fmt.Printf("Error: Parsing request body (%s)\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Internal Server Error"})
		return
	}

	// Validate the submission
	if errs := h.app.ValidateSubmissionWithForm(jsonMap, form); len(errs) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "There was an error validating the submission"})
		return
	}

	// Validate the reCAPTCHA token
	if err := h.app.ValidateRecaptchaToken(jsonMap["recaptcha_token"], c.ClientIP()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the submission
	if submission, err := h.app.CreateSubmission(c.Param("slug"), c.ClientIP(), c.GetHeader("User-Agent"), jsonMap); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error creating the submission"})
		return
	} else {
		c.JSON(http.StatusCreated, submission)
	}
}

func (h *httpLayer) getSubmissionsApi(c *gin.Context) {

}
