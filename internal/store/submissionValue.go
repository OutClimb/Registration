package store

type SubmissionValue struct {
	ID           uint `gorm:"primaryKey"`
	SubmissionID uint
	FormValueID  uint
	Value        string
}

func (s *storeLayer) CreateSubmissionValue(submissionId uint, formValueId uint, value string) (SubmissionValue, error) {
	submissionValue := SubmissionValue{
		SubmissionID: submissionId,
		FormValueID:  formValueId,
		Value:        value,
	}

	if result := s.db.Create(submissionValue); result.Error != nil {
		return SubmissionValue{}, result.Error
	}

	return submissionValue, nil
}
