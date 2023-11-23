package mocks

import (
	"gistapp.ck89.net/internal/dblayer"
	"time"
)

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

func (u *UserLayer) Fetch(id int) (dblayer.User, error) {
	if id == 1 {
		u := dblayer.User{
			ID:      1,
			Name:    "cksre",
			Email:   "cksre@gmail.com",
			Created: time.Now(),
		}
		return u, nil
	}
	return dblayer.User{}, dblayer.ErrNoRecord
}
