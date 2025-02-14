package app

import (
	"errors"
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
