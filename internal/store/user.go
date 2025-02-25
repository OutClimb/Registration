package store

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"uniqueIndex;not null;size:255"`
	Password  string `gorm:"not null;size:64"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"not null"`
	Token     string
	IPAddress string
	Disabled  bool
}
