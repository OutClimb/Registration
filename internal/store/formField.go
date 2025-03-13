package store

import "database/sql"

type FormField struct {
	ID         uint `gorm:"primaryKey"`
	FormID     uint
	Name       string `gorm:"not null"`
	Slug       string `gorm:"not null;size:255"`
	Type       string `gorm:"not null;size:32"`
	Metadata   sql.NullString
	Required   bool
	Validation sql.NullString
	Order      uint `gorm:"not null;default:0"`
}

func (s *storeLayer) GetFormFields(formID uint) (*[]FormField, error) {
	fields := []FormField{}

	if result := s.db.Where("form_id = ?", formID).Find(&fields); result.Error != nil {
		return &[]FormField{}, result.Error
	}

	return &fields, nil
}
