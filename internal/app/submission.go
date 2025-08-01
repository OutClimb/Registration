package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/OutClimb/Registration/internal/store"
)

type SubmissionsInternal map[string]string

type SubmissionInternal struct {
	SubmittedOn time.Time
}

type EmailSubmissionData struct {
	Form   *store.Form
	Fields *[]store.FormField
	Values *map[string]string
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
		emailData := EmailSubmissionData{
			Form:   form,
			Fields: fields,
			Values: &values,
		}

		if len(form.EmailTo) != 0 {
			a.sendEmail("noreply@outclimb.gay", form.EmailTo, "Form Submission - "+form.Name, "internal", emailData)
		}

		if form.EmailFormFieldID != 0 && form.EmailSubject != "" && form.EmailTemplate != "" {
			emailTo := ""
			for _, field := range *fields {
				if field.ID == form.EmailFormFieldID {
					emailTo = values[field.Slug]
				}
			}

			if emailTo != "" {
				a.sendEmail("noreply@outclimb.gay", emailTo, form.EmailSubject, form.EmailTemplate, emailData)
			}
		}

		submissionInternal := SubmissionInternal{}
		submissionInternal.Internalize(submission)

		return &submissionInternal, nil
	}
}

func (a *appLayer) DeleteSubmission(id uint) error {
	if err := a.store.DeleteSubmission(id); err != nil {
		return errors.New("Unable to delete submission")
	}

	return nil
}

func (a *appLayer) GetSubmissionsForForm(slug string, userId uint) (*[]SubmissionsInternal, error) {
	// Get the user.
	user, err := a.store.GetUser(userId)
	if err != nil {
		return nil, errors.New("User not found")
	}

	// Check if user is admin or viewer.
	if user.Role != "admin" && user.Role != "viewer" {
		return nil, errors.New("Unauthorized")
	}

	// Get the form.
	form, err := a.store.GetForm(slug)
	if err != nil {
		return nil, errors.New("Form not found")
	}

	// Check if user is allowed to view form.
	if user.Role == "viewer" {
		allowed := false
		for _, u := range form.ViewableBy {
			if u.ID == userId {
				allowed = true
				break
			}
		}

		if !allowed {
			return nil, errors.New("Unauthorized")
		}
	}

	// Get the submissions.
	submissions, err := a.store.GetSubmissionsForForm(form.ID)
	if err != nil {
		return nil, errors.New("Error getting submissions")
	}

	// Get the form fields.
	fields, err := a.store.GetFormFields(form.ID)
	if err != nil {
		return nil, errors.New("Error getting form fields")
	}

	// Transform the form fields into a map.
	fieldMap := make(map[uint]store.FormField, len(*fields))
	for _, field := range *fields {
		fieldMap[field.ID] = field
	}

	// Get the submission values.
	result := make([]SubmissionsInternal, len(*submissions))
	for i, submission := range *submissions {
		values, err := a.store.GetSubmissionValues(submission.ID)
		if err != nil {
			return nil, errors.New("Error getting submission values")
		}

		submissionMap := make(SubmissionsInternal)
		submissionMap["id"] = strconv.FormatUint(uint64(submission.ID), 10)
		submissionMap["submitted_on"] = submission.SubmittedOn.Format(time.UnixDate)
		submissionMap["ip_address"] = submission.IPAddress
		submissionMap["user_agent"] = submission.UserAgent
		for _, value := range *values {
			field := fieldMap[value.FormFieldID]
			submissionMap[field.Slug] = value.Value
		}

		result[i] = submissionMap
	}

	return &result, nil
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

func (a *appLayer) ValidateSubmissionWithForm(submission map[string]string, form *FormInternal) []error {
	// Make sure reCAPTCHA token is present
	if _, ok := submission["recaptcha_token"]; !ok {
		return []error{
			errors.New("Missing required field: recaptcha_token"),
		}
	}

	errs := []error{}
	for _, field := range form.Fields {
		// Validate required fields.
		if field.Required {
			if value, ok := submission[field.Slug]; !ok || len(strings.TrimSpace(value)) == 0 {
				errs = append(errs, errors.New("Missing required field: "+field.Name))
			}
		}

		// Validate fields with validation.
		if field.Validation != "" && len(strings.TrimSpace(submission[field.Slug])) != 0 {
			if matched, _ := regexp.MatchString(field.Validation, submission[field.Slug]); !matched {
				errs = append(errs, errors.New("Field "+field.Name+" does not match validation"))
			}
		}

		// Validate checkbox fields.
		if field.Type == "checkboxes" && submission[field.Slug] != "" {
			options := (*form.Fields[field.Slug].Metadata).(map[string]interface{})
			selectedOptions := strings.Split(submission[field.Slug], ", ")
			for _, selectedOption := range selectedOptions {
				if _, ok := options[selectedOption]; !ok {
					errs = append(errs, errors.New("Invalid option selected for "+field.Name))
				}
			}
		}

		// Validate select fields.
		if (field.Type == "radios" || field.Type == "select") && submission[field.Slug] != "" {
			options := (*form.Fields[field.Slug].Metadata).(map[string]interface{})
			if _, ok := options[submission[field.Slug]]; !ok {
				errs = append(errs, errors.New("Invalid option selected for "+field.Name))
			}
		}

		// Validate boolean fields.
		if field.Type == "bool" && submission[field.Slug] != "" {
			lowerValue := strings.ToLower(submission[field.Slug])
			possibleValues := map[string]bool{"true": true, "false": true, "1": true, "0": true, "yes": true, "no": true}
			if _, ok := possibleValues[lowerValue]; !ok {
				errs = append(errs, errors.New("Invalid value for "+field.Name))
			}
		}
	}

	return errs
}
