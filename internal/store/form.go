package store

import (
	"time"
)

type Form struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"not null"`
	Slug           string `gorm:"uniqueIndex;not null;size:255"`
	Template       string `gorm:"not null;size:245"`
	OpensOn        time.Time
	ClosesOn       time.Time
	MaxSubmissions uint
	FormFields     []FormField
	Submissions    []Submission
}

func (s *storeLayer) GetForm(slug string) (Form, error) {
	form := Form{}

	if result := s.db.Where("slug = ?", slug).First(&form); result.Error != nil {
		return Form{}, result.Error
	}

	return form, nil
}

func (s *storeLayer) GetFormExists(slug string) bool {
	var count int64 = 0
	s.db.Model(&Form{}).Where("slug = ?", slug).Count(&count)

	return count == 1
}
