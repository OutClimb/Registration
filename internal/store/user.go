package store

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username             string `gorm:"uniqueIndex;not null;size:255"`
	Password             string `gorm:"not null;size:64"`
	Role                 string `gorm:"not null;size:255"`
	Name                 string `gorm:"not null"`
	Email                string `gorm:"not null"`
	Disabled             bool
	RequirePasswordReset bool `gorm:"not null;default:0"`
}

func (s *storeLayer) CreateUser(username, password, name, email, role string) (*User, error) {
	user := User{
		Username: username,
		Password: password,
		Name:     name,
		Email:    email,
		Role:     role,
	}

	if result := s.db.Create(&user); result.Error != nil {
		return &User{}, result.Error
	}

	return &user, nil
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

func (s *storeLayer) UpdatePassword(id uint, password string) error {
	user, err := s.GetUser(id)
	if err != nil {
		return err
	}

	user.Password = password
	user.RequirePasswordReset = false

	if result := s.db.Save(user); result.Error != nil {
		return result.Error
	}

	return nil
}
