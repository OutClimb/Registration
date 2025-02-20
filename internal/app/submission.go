package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/OutClimb/Registration/internal/store"
)

type SubmissionInternal struct {
	SubmittedOn time.Time
}

func (s *SubmissionInternal) Internalize(field *store.Submission) {
	s.SubmittedOn = field.SubmittedOn
}

type SiteVerifyResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func (a *appLayer) CreateSubmission(slug string, ipAddress string, userAgent string, values map[string]string) (*SubmissionInternal, error) {
	if form, err := a.store.GetForm(slug); err != nil {
		return &SubmissionInternal{}, err
	} else if fields, err := a.store.GetFormFields(form.ID); err != nil {
		return &SubmissionInternal{}, err
	} else if submission, err := a.store.CreateSubmission(form, fields, ipAddress, userAgent, values); err != nil {
		return &SubmissionInternal{}, err
	} else {
		submissionInternal := SubmissionInternal{}
		submissionInternal.Internalize(submission)

		return &submissionInternal, nil
	}
}

func (a *appLayer) ValidateRecaptchaToken(token string, clientIp string) error {
	const siteVerifyURL = "https://www.google.com/recaptcha/api/siteverify"
	secretKey, _ := os.LookupEnv("RECAPTCHA_SECRET_KEY")

	req, err := http.NewRequest(http.MethodPost, siteVerifyURL, nil)
	if err != nil {
		fmt.Printf("Error: Issue creating reCAPTCHA request (%s)\n", err.Error())
		return errors.New("Internal reCAPTCHA error")
	}

	// Add necessary request parameters.
	q := req.URL.Query()
	q.Add("secret", secretKey)
	q.Add("response", token)
	q.Add("remoteip", clientIp)
	req.URL.RawQuery = q.Encode()

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error: Issue making reCAPTCHA request (%s)\n", err.Error())
		return errors.New("Internal reCAPTCHA error")
	}
	defer resp.Body.Close()

	// Decode response.
	var body SiteVerifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		fmt.Printf("Error: Issue decoding reCAPTCHA response (%s)\n", err.Error())
		return errors.New("Internal reCAPTCHA error")
	}

	// Check recaptcha verification success.
	if !body.Success {
		fmt.Println("Error: Unsuccessful reCAPTCHA verify request")
		return errors.New("reCAPTCHA unsuccessful")
	}

	// Check response score.
	if body.Score < 0.5 {
		fmt.Printf("Error: reCAPTCHA score lower that 0.5 (%f)\n", body.Score)
		return errors.New("reCAPTCHA unsuccessful")
	}

	// Check response action.
	if body.Action != "submit" {
		fmt.Printf("Error: Mismatched reCAPTCHA action (%s)\n", body.Action)
		return errors.New("reCAPTCHA unsuccessful")
	}

	return nil
}

func (a *appLayer) ValidateSubmissionWithForm(submission map[string]string, form *FormInternal) error {
	for _, field := range form.Fields {
		if field.Required {
			if value, ok := submission[field.Slug]; !ok || len(strings.TrimSpace(value)) == 0 {
				return errors.New("Missing required field: " + field.Name)
			}
		}

		if field.Validation != "" {
			if matched, _ := regexp.MatchString(field.Validation, submission[field.Slug]); !matched {
				return errors.New("Field " + field.Name + " does not match validation")
			}
		}
	}

	return nil
}
