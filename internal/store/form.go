package store

import (
	"time"
)

type Form struct {
	ID                 uint   `gorm:"primaryKey"`
	Name               string `gorm:"not null"`
	Slug               string `gorm:"uniqueIndex;not null;size:255"`
	Template           string `gorm:"not null;size:245"`
	OpensOn            *time.Time
	ClosesOn           *time.Time
	MaxSubmissions     uint
	ViewableBy         []User `gorm:"many2many:form_viewable_users;"`
	NotOpenMessage     *string
	ClosedMessage      *string
	SuccessMessage     *string
	EmailFormFieldSlug string
	EmailTo            string
	EmailSubject       string
	EmailTemplate      string
}

func (s *storeLayer) CreateForm(name, slug, template string, opensOn, closesOn *time.Time, maxSubmissions uint, notOpenMessage, closedMessage, successMessage *string, emailFormFieldSlug, emailTo, emailSubject, emailTemplate string) (*Form, error) {
	form := Form{
		Name:               name,
		Slug:               slug,
		Template:           template,
		OpensOn:            opensOn,
		ClosesOn:           closesOn,
		MaxSubmissions:     maxSubmissions,
		NotOpenMessage:     notOpenMessage,
		ClosedMessage:      closedMessage,
		SuccessMessage:     successMessage,
		EmailFormFieldSlug: emailFormFieldSlug,
		EmailTo:            emailTo,
		EmailSubject:       emailSubject,
		EmailTemplate:      emailTemplate,
	}

	if result := s.db.Create(&form); result.Error != nil {
		return nil, result.Error
	}

	return &form, nil
}

func (s *storeLayer) GetAllForms() (*[]Form, error) {
	forms := []Form{}

	if result := s.db.Find(&forms); result.Error != nil {
		return &[]Form{}, result.Error
	}

	return &forms, nil
}

func (s *storeLayer) GetForm(slug string) (*Form, error) {
	form := Form{}

	if result := s.db.Model(&Form{}).Preload("ViewableBy").Where("slug = ?", slug).First(&form); result.Error != nil {
		return &Form{}, result.Error
	}

	return &form, nil
}

func (s *storeLayer) GetFormsForUser(userId uint) (*[]Form, error) {
	forms := []Form{}

	if result := s.db.Joins("LEFT JOIN form_viewable_users ON id = form_id").Where("user_id = ?", userId).Find(&forms); result.Error != nil {
		return &[]Form{}, result.Error
	}

	return &forms, nil
}
