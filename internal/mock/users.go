package mock

import (
	"github.com/huytran2000-hcmus/snippetbox/internal/models"
)

type StubUsers struct{}

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
