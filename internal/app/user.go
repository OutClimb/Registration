package app

import (
	"errors"
	"os"
	"strconv"
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
	Role      string
	Name      string
	Email     string
	Disabled  bool
}

func (u *UserInternal) Internalize(user *store.User) {
	u.ID = user.ID
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	u.Username = user.Username
	u.Role = user.Role
	u.Name = user.Name
	u.Email = user.Email
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

func (a *appLayer) CheckRole(userRole string, requiredRole string) bool {
	roleMap := map[string]int{
		"admin":  3,
		"viewer": 2,
		"user":   1,
	}

	if roleMap[userRole] >= roleMap[requiredRole] {
		return true
	}

	return false
}

func (a *appLayer) CreateToken(user *UserInternal, clientIp string) (string, error) {
	// Get the token lifespan
	tokenLifespan, err := strconv.Atoi(os.Getenv("TOKEN_LIFESPAN"))
	if err != nil {
		return "", errors.New("Failed to get token lifespan")
	}

	// Create the Claims
	claims := jwt.MapClaims{}
	claims["user_id"] = user.ID
	claims["ip_address"] = clientIp
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	claims["iss"] = "api"

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))); err != nil {
		return "", errors.New("Failed to sign token")
	} else {
		return signedToken, nil
	}
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
