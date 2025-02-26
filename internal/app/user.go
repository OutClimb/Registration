package app

import (
	"errors"
	"os"
	"time"

	"github.com/OutClimb/Registration/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserInternal struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Username  string
	Name      string
	Email     string
	IPAddress string
	Disabled  bool
}

func (u *UserInternal) Internalize(user *store.User) {
	u.ID = user.ID
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	u.Username = user.Username
	u.Name = user.Name
	u.Email = user.Email
	u.IPAddress = user.IPAddress
	u.Disabled = user.Disabled

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

func (a *appLayer) CreateToken(user *UserInternal, clientIp string) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{}
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["iss"] = "api"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))); err != nil {
		return "", errors.New("Failed to sign token")
	} else if err := a.store.UpdateUserWithToken(user.ID, signedToken, clientIp); err != nil {
		return "", errors.New("Failed to update user token")
	} else {
		return signedToken, nil
	}
}

func (a *appLayer) ValidateToken(userId uint, clientIp string) error {
	if user, err := a.store.GetUser(userId); err != nil {
		return errors.New("User not found")
	} else {
		if user.Disabled {
			return errors.New("User is disabled")
		}

		if user.IPAddress != clientIp {
			return errors.New("Wrong Location")
		}
	}

	return nil
}
