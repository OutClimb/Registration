package store

import (
	"time"

	"gorm.io/gorm"
)

type Submission struct {
	ID          uint `gorm:"primaryKey"`
	FormID      uint
	SubmittedOn time.Time
	IPAddress   string
	UserAgent   string
}

func (s *storeLayer) CreateSubmission(form *Form, fields *[]FormField, ipAddress string, userAgent string, values map[string]string) (*Submission, error) {
	submission := Submission{
		FormID:      form.ID,
		SubmittedOn: time.Now(),
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}

	s.db.Transaction(func(tx *gorm.DB) error {
		if result := s.db.Create(&submission); result.Error != nil {
			return result.Error
		}

		for _, field := range *fields {
			if _, err := s.CreateSubmissionValue(submission.ID, field.ID, values[field.Slug]); err != nil {
				return err
			}
		}

		return nil
	})

	return &submission, nil
}

func (s *storeLayer) GetNumberOfSubmissions(formID uint) (int64, error) {
	var count int64

	if result := s.db.Model(&Submission{}).Where("form_id = ?", formID).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (s *storeLayer) GetSubmissionsForForm(formID uint) (*[]Submission, error) {
	submissions := []Submission{}
	if result := s.db.Where("form_id = ?", formID).Find(&submissions); result.Error != nil {
		return &[]Submission{}, result.Error
	}

	return &submissions, nil
}
