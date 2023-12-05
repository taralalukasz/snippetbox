package mysql

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"tarala/snippetbox/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, hashedPassword)

	//checking, if the error code is 1062 (mysql can't insert) and the reason of it.
	//If email constraint is violated - then return error
	if err != nil {
		if sqlError, ok := err.(*mysql.MySQLError); ok {
			if sqlError.Number == 1062 && strings.Contains(sqlError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
	}

	return err
}

func (m UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// We'll use the Get method to fetch details for a specific user based
// on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
