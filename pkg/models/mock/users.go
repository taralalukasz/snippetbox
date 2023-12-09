package mock

import (
	"time"

	"tarala/snippetbox/pkg/models"
)

type UserModel struct{}

var mockUser = &models.User{
	ID:      1,
	Name:    "lukasz",
	Email:   "lukasz@gmail.com",
	Created: time.Now(),
}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "wrong@gmail.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	switch email {
	case "lukasz@gmail.com":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	} 
	return 1, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1 :
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord 
	}
}
