package store

import (
	"time"
)

type Submission struct {
	ID               uint `gorm:"primaryKey"`
	FormID           uint
	SubmittedOn      time.Time
	IPAddress        string
	UserAgent        string
	SubmissionValues []SubmissionValue
}

func (s *storeLayer) CreateSubmission(formId uint, ipAddress string, userAgent string, values map[uint]string) (Submission, error) {
	submission := Submission{
		FormID:      formId,
		SubmittedOn: time.Now(),
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}

	if result := s.db.Create(submission); result.Error != nil {
		return Submission{}, result.Error
	}

	submission.SubmissionValues = make([]SubmissionValue, len(values))
	index := 0
	for fieldId, value := range values {
		if submissionValue, error := s.CreateSubmissionValue(submission.ID, fieldId, value); error == nil {
			submission.SubmissionValues[index] = submissionValue
			index++
		} else {
			// TODO: Rollback submission
		}
	}

	return submission, nil
}
