package models

import (
	"errors"
)

var (

	// error if no snippet is found for specified record id
	ErrNoRecord = errors.New("models: no matching record found")

	// error if a user tries to login with an incorrect email address or password
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	//error if a user tries to signup with an email address that is already in use
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
