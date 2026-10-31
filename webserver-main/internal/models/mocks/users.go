package mocks

import (
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
