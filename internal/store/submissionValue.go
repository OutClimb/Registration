package store

type SubmissionValue struct {
	ID           uint `gorm:"primaryKey"`
	SubmissionID uint
	FormFieldID  uint
	Value        string
}

func (s *storeLayer) CreateSubmissionValue(submissionId uint, formFieldId uint, value string) (*SubmissionValue, error) {
	if value == "" {
		return &SubmissionValue{}, nil
	}

	submissionValue := SubmissionValue{
		SubmissionID: submissionId,
		FormFieldID:  formFieldId,
		Value:        value,
	}

	if result := s.db.Create(&submissionValue); result.Error != nil {
		return &SubmissionValue{}, result.Error
	}

	return &submissionValue, nil
}

func (s *storeLayer) GetSubmissionValues(submissionId uint) (*[]SubmissionValue, error) {
	submissionValues := []SubmissionValue{}
	if result := s.db.Where("submission_id = ?", submissionId).Find(&submissionValues); result.Error != nil {
		return &[]SubmissionValue{}, result.Error
	}

	return &submissionValues, nil
}
