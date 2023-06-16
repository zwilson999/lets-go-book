package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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

	// use cost of 12 which is a sensible number
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `
		INSERT INTO
			users (name, email, hashed_password, created)
		VALUES(
			?, ?, ?, UTC_TIMESTAMP()
		)
	`
	// use the Exec()m ethod to insert the user details and hashed password into the users table
	_, err = um.DB.Exec(stmt, name, email, string(hashedPass))
	if err != nil {
		// if this returns an error, we use errors.As() to check if it has a specific mysql error type
		// if it does, the error will be assigned to the mySQLError variable. We can check whether
		// or not the error relates to our users_uc_email key by checking if the error code equals 1062 and the contents of the error message string
		// if it does we will return an ErrDuplicateEmail error
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// method to authenticate to verify whether a user exists with the provided email
// and password. this will return the relevant user ID if they do.
func (um *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email =?"

	err := um.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// check whether the hashed password and plain-text password provided, match.
	// if they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// otherwise the id is correct
	return id, nil
}

// method to check if a user exists with a specific ID.
func (um *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
