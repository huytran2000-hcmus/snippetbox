package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserRepository struct {
	DB *sql.DB
}

func (rep *UserRepository) Insert(name string, email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return ErrPasswordTooLong
		}

		return err
	}

	stmt := "INSERT INTO users (name, email, hashed_password, created) VALUES ($1, $2, $3, NOW())"
	_, err = rep.DB.Exec(stmt, name, email, hashedPassword)
	if err != nil {
		var postgresErr *pq.Error
		if errors.As(err, &postgresErr); postgresErr != nil {
			if postgresErr.Code == "23505" && postgresErr.Constraint == "users_uc_email" {
				return ErrDuplicateEmail
			}
		}

		return fmt.Errorf("models: error when inserting a user: %s", err)
	}

	return nil
}

func (rep *UserRepository) Authenticate(email string, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users where email = $1"
	err := rep.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}

		return 0, fmt.Errorf("models: error when selecting a user: %s", err)
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}

		return 0, err
	}

	return id, nil
}

func (rep *UserRepository) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = $1)"

	err := rep.DB.QueryRow(stmt, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrNoRecord
		}

		return false, fmt.Errorf("models: errors when selecting a user")
	}

	return exists, nil
}
