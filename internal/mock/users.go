package mock

import (
	"time"

	"github.com/huytran2000-hcmus/snippetbox/internal/models"
)

var mockUser = &models.User{
	ID:             1,
	Name:           "alice",
	Email:          "dupe@example.com",
	HashedPassword: []byte{},
	Created:        time.Date(2023, time.May, 10, 20, 0, 0, 0, time.UTC),
}

type StubUsers struct{}

func (s *StubUsers) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (s *StubUsers) Insert(name string, email string, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (s *StubUsers) Authenticate(email string, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (s *StubUsers) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
