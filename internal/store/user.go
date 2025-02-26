package store

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null;size:255"`
	Password string `gorm:"not null;size:64"`
	Role     string `gorm:"not null;size:255"`
	Name     string `gorm:"not null"`
	Email    string `gorm:"not null"`
	Disabled bool
}

func (s *storeLayer) GetUser(id uint) (*User, error) {
	user := User{}

	if result := s.db.Where("id = ?", id).First(&user); result.Error != nil {
		return &User{}, result.Error
	}

	return &user, nil
}

func (s *storeLayer) GetUserWithUsername(username string) (*User, error) {
	user := User{}

	if result := s.db.Where("username = ?", username).First(&user); result.Error != nil {
		return &User{}, result.Error
	}

	return &user, nil
}
