package app

import (
	"errors"
)

func (a *appLayer) ValidateUserToken(userId uint, clientIp string) error {
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
