package models

import (
	"database/sql"
	"time"
)

// define a User type that has types that align with our database column types
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

// method to insert new record into our users table
func (um *UserModel) Insert(name, email, password string) error {
	return nil
}

// method to authenticate to verify whether a user exists with the provided email
// and password. this will return the relevant user ID if they do.
func (um *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// method to check if a user exists with a specific ID.
func (um *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
