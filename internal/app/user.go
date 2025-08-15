package app

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/OutClimb/Registration/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type UserInternal struct {
	ID                   uint
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time
	Username             string
	Password             string
	Role                 string
	Name                 string
	Email                string
	Disabled             bool
	RequirePasswordReset bool
}

func (u *UserInternal) Internalize(user *store.User) {
	u.ID = user.ID
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	u.Username = user.Username
	u.Password = user.Password
	u.Role = user.Role
	u.Name = user.Name
	u.Email = user.Email
	u.Disabled = user.Disabled
	u.RequirePasswordReset = user.RequirePasswordReset

	if user.DeletedAt.Valid {
		u.DeletedAt = &user.DeletedAt.Time
	}
}

func (a *appLayer) AuthenticateUser(username string, password string) (*UserInternal, error) {
	if user, err := a.store.GetUserWithUsername(username); err != nil {
		return &UserInternal{}, err
	} else if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return &UserInternal{}, errors.New("Invalid password")
	} else {
		userInternal := UserInternal{}
		userInternal.Internalize(user)

		return &userInternal, nil
	}
}

func (a *appLayer) CheckRole(userRole string, requiredRole string) bool {
	roleMap := map[string]int{
		"admin":  3,
		"viewer": 2,
		"user":   1,
		"reset":  0,
	}

	if roleMap[userRole] >= roleMap[requiredRole] {
		return true
	}

	return false
}

func (a *appLayer) GetUser(userId uint) (*UserInternal, error) {
	if user, err := a.store.GetUser(userId); err != nil {
		return &UserInternal{}, err
	} else {
		userInternal := UserInternal{}
		userInternal.Internalize(user)

		return &userInternal, nil
	}
}

func (a *appLayer) UpdatePassword(user *UserInternal, password string) error {
	cost, err := strconv.Atoi(os.Getenv("PASSWORD_COST"))
	if err != nil {
		return errors.New("Unknown cost")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return errors.New("Failed to hash password")
	}

	if err := a.store.UpdatePassword(user.ID, string(hashedPassword)); err != nil {
		return errors.New("Failed to update password")
	}

	return nil
}

func (a *appLayer) ValidateUser(userId uint) error {
	if user, err := a.store.GetUser(userId); err != nil {
		return errors.New("User not found")
	} else {
		if user.Disabled {
			return errors.New("User is disabled")
		}
	}

	return nil
}

func (a *appLayer) ValidatePassword(user *UserInternal, password string) error {
	if len(password) < 16 {
		return errors.New("Password must be at least 16 characters")
	} else if len(password) > 72 {
		return errors.New("Password must be at most 72 characters")
	} else if symbolMatched, _ := regexp.MatchString("[^a-zA-Z0-9]", password); !symbolMatched {
		return errors.New("Password must contain a symbol")
	} else if numberMatched, _ := regexp.MatchString("[0-9]", password); !numberMatched {
		return errors.New("Password must contain a number")
	} else if upperMatched, _ := regexp.MatchString("[A-Z]", password); !upperMatched {
		return errors.New("Password must contain an uppercase letter")
	} else if lowerMatched, _ := regexp.MatchString("[a-z]", password); !lowerMatched {
		return errors.New("Password must contain a lowercase letter")
	} else if strings.Contains(password, user.Username) {
		return errors.New("Password must not contain the username")
	} else if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
		return errors.New("Password must be different from the current password")
	}

	return nil
}
