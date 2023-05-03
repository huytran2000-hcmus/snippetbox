package models

import (
	"database/sql"
	"errors"
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

func (rst *UserRepository) Insert(name string, email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO users (name, email, hashed_password, created) VALUES ($1, $2, $3, NOW())"
	_, err = rst.DB.Exec(stmt, name, email, hashedPassword)
	if err != nil {
		var postgresErr *pq.Error
		if errors.As(err, &postgresErr); postgresErr != nil {
			if postgresErr.Code == "23505" && postgresErr.Constraint == "users_uc_email" {
				return ErrDuplicateEmail
			}
		}

		return err
	}

	return nil
}
