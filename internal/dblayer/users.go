package dblayer

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserLayer struct {
	MysqlDB *sql.DB
}

type UserLayerIface interface {
	Add(name, email, password string) error
	Authn(email, password string) (int, error)
	CheckExists(id int) (bool, error)
	Fetch(id int) (User, error)
}

func (m *UserLayer) Fetch(id int) (User, error) {
	var user User
	stmt := `SELECT id, name, email, created FROM users WHERE id = ?`

	err := m.MysqlDB.QueryRow(stmt, id).Scan(&user.ID, &user.Name, &user.Email, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord
		} else {
			return User{}, err
		}
	}
	return user, nil
}

func (m *UserLayer) Add(name, email, password string) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	statement := `INSERT INTO users (name, email, hashed_password, created) VALUES (?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.MysqlDB.Exec(statement, name, email, hashedPwd)
	if err != nil {
		var dbError *mysql.MySQLError
		if errors.As(err, &dbError) {
			if dbError.Number == 1062 && strings.Contains(dbError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserLayer) Authn(email, password string) (int, error) {
	var id int
	var hashedPwd []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
	err := m.MysqlDB.QueryRow(stmt, email).Scan(&id, &hashedPwd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPwd, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserLayer) CheckExists(id int) (bool, error) {
	var isPresent bool
	stmt := "SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)"
	err := m.MysqlDB.QueryRow(stmt, id).Scan(&isPresent)
	return isPresent, err
}
