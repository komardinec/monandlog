package mocks

import (
	"time"

	"debutwebapp.komardinec/internal/models"
)

type UserModel struct{}

func (u *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@dupe.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	if email == "example@example.com" && password == "P@ssw0rd" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (u *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (u *UserModel) Get(id int) (models.User, error) {
	switch id {
	case 1:
		user := models.User{
			ID:      1,
			Name:    "example",
			Email:   "example@example.com",
			Created: time.Now(),
		}
		return user, nil
	default:
		return models.User{}, models.ErrNoRecord
	}
}

func (u *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	return nil
}
