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

type Users interface {
	Get(id int) (*User, error)
	Insert(name string, email string, password string) error
	Authenticate(email string, password string) (int, error)
	Exists(id int) (bool, error)
	PasswordUpdate(id int, currentPassword string, newPassword string) error
}

type UserDB struct {
	DB *sql.DB
}

func (db *UserDB) Get(id int) (*User, error) {
	stmt := "SELECT name, email, created FROM users WHERE id = $1"

	var user User
	err := db.DB.QueryRow(stmt, id).Scan(&user.Name, &user.Email, &user.Created)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}

		return nil, fmt.Errorf("models: select a user: %s", err)
	}

	return &user, nil
}

const passwordHashingCost = 12

func (db *UserDB) Insert(name string, email string, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO users (name, email, hashed_password, created) VALUES ($1, $2, $3, NOW())"
	_, err = db.DB.Exec(stmt, name, email, hashedPassword)
	if err != nil {
		var postgresErr *pq.Error
		if errors.As(err, &postgresErr); postgresErr != nil {
			if postgresErr.Code == "23505" && postgresErr.Constraint == "users_uc_email" {
				return ErrDuplicateEmail
			}
		}

		return fmt.Errorf("models: insert a user: %s", err)
	}

	return nil
}

func (db *UserDB) Authenticate(email string, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users where email = $1"
	err := db.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}

		return 0, fmt.Errorf("models: select a user: %s", err)
	}

	err = compareHashedPassword(hashedPassword, []byte(password))
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *UserDB) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = $1)"

	err := db.DB.QueryRow(stmt, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrNoRecord
		}

		return false, fmt.Errorf("models: select an existing user: %s", err)
	}

	return exists, nil
}

func (db *UserDB) PasswordUpdate(id int, currentPassword string, newPassword string) error {
	stmt := "SELECT hashed_password FROM users WHERE id = $1"

	var hashedPassword []byte
	err := db.DB.QueryRow(stmt, id).Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}

		return err
	}

	err = compareHashedPassword(hashedPassword, []byte(currentPassword))
	if err != nil {
		return err
	}

	hashedPassword, err = hashPassword(newPassword)
	if err != nil {
		return err
	}

	stmt = "UPDATE users SET hashed_password = $1 WHERE id = $2"
	_, err = db.DB.Exec(stmt, hashedPassword, id)
	if err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashingCost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return nil, ErrPasswordTooLong
		}

		return nil, fmt.Errorf("models: hash password: %s", err)
	}

	return hashedPassword, nil
}

func compareHashedPassword(hashedPassword, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}

		return fmt.Errorf("models: compare password to hashed password: %s", err)
	}

	return nil
}
