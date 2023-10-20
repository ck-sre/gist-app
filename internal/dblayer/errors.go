package dblayer

import "errors"

var (
	ErrNoRecord           = errors.New("dblayer: no matching record found")
	ErrInvalidCredentials = errors.New("dblayer: invalid credentials")
	ErrDuplicateEmail     = errors.New("dblayer: duplicate email")
)
