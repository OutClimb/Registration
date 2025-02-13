package store

type FormField struct {
	ID               uint `gorm:"primaryKey"`
	FormID           uint
	Name             string `gorm:"not null"`
	Slug             string `gorm:"not null;size:255"`
	Type             string `gorm:"not null;size:32"`
	Metadata         string
	Required         bool
	Validation       string
	SubmissionValues []SubmissionValue
}

func (s *storeLayer) GetFormFields(formID uint) ([]FormField, error) {
	fields := []FormField{}

	if result := s.db.Where("form_id = ?", formID).Find(&fields); result.Error != nil {
		return []FormField{}, result.Error
	}

	return fields, nil
}

func (s *storeLayer) GetFormField(formID uint, slug string) (FormField, error) {
	field := FormField{}

	if result := s.db.Where("form_id = ? AND slug = ?", formID, slug).First(&field); result.Error != nil {
		return FormField{}, result.Error
	}

	return field, nil
}
