package mysql

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"snippetbox/pkg/models"
	"strings"
)

type UserModel struct {
	DB *sql.DB
}

// Insert a single user to the database, the password is hashed using bcrypt
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) {
			if mysqlError.Number == 1062 && strings.Contains(mysqlError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil

}

// Authenticate is a helper method to validate whether a given password is correct. It comapres a hash of the supplied
// password with the password hash stored in the databse
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

// Get retrieves a User given an int id.
func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}

	stmt := `SELECT id, name, email, created, active FROM users WHERE id =?`
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return u, nil
}
