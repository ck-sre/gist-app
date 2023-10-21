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

func (m *UserLayer) Auth(email, password string) (int, error) {
	return 0, nil
}

func (m *UserLayer) CheckExists(id int) (bool, error) {
	return false, nil
}