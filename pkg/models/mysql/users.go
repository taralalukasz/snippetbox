package mysql

import (
	"database/sql"
	"tarala/snippetbox/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(name, email, password string) error {
	return nil
}

func (m UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// We'll use the Get method to fetch details for a specific user based
// on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
