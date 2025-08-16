package http

import (
	"github.com/OutClimb/Registration/internal/app"
)

type userPublic struct {
	Username             string `json:"username"`
	Role                 string `json:"role"`
	Name                 string `json:"name"`
	Email                string `json:"email"`
	RequirePasswordReset bool   `json:"requirePasswordReset"`
}

func (u *userPublic) Publicize(user *app.UserInternal) {
	u.Username = user.Username
	u.Role = user.Role
	u.Name = user.Name
	u.Email = user.Email
	u.RequirePasswordReset = user.RequirePasswordReset
}
