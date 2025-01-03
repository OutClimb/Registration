package store

type Event struct {
	Name string `gorm:"not null"`
	Slug string `gorm:"uniqueIndex;not null;size:512"`
}

func (s *storeLayer) CreateEvent(name string, slug string) (Event, error) {
	event := Event{
		Name: name,
		Slug: slug,
	}

	if result := s.db.Create(event); result.Error != nil {
		return Event{}, result.Error
	}

	return event, nil
}

func (s *storeLayer) GetEvent(slug string) (Event, error) {
	event := Event{}

	if result := s.db.Where("slug = ?", slug).First(&event); result.Error != nil {
		return Event{}, result.Error
	}

	return event, nil
}
