package store

type FormField struct {
	ID         uint `gorm:"primaryKey"`
	FormID     uint
	Name       string `gorm:"not null"`
	Slug       string `gorm:"not null;size:255"`
	Type       string `gorm:"not null;size:32"`
	Metadata   *string
	Required   bool
	Validation *string
	Order      uint `gorm:"not null;default:0"`
}

func (s *storeLayer) CreateFormField(formID uint, name, slug, fieldType string, metadata *string, required bool, validation *string, order uint) (*FormField, error) {
	field := FormField{
		FormID:     formID,
		Name:       name,
		Slug:       slug,
		Type:       fieldType,
		Metadata:   metadata,
		Required:   required,
		Validation: validation,
		Order:      order,
	}

	if result := s.db.Create(&field); result.Error != nil {
		return nil, result.Error
	}

	return &field, nil
}

func (s *storeLayer) GetFormFields(formID uint) (*[]FormField, error) {
	fields := []FormField{}

	if result := s.db.Where("form_id = ?", formID).Find(&fields); result.Error != nil {
		return &[]FormField{}, result.Error
	}

	return &fields, nil
}
