package mocks

import "gistapp.ck89.net/internal/dblayer"

type UserLayer struct{}

func (u *UserLayer) Add(name, email, password string) error {
	switch email {
	case "cksre@gmail.com":
		return dblayer.ErrDuplicateEmail
	default:
		return nil

	}
}

func (u *UserLayer) Authn(email, password string) (int, error) {
	if email == "cksre@gmail.com" && password == "cksregmailcom" {
		return 1, nil
	}
	return 0, dblayer.ErrInvalidCredentials
}

func (u *UserLayer) CheckExists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
