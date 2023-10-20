package dblayer

import (
	"database/sql"
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
	return nil
}

func (m *UserLayer) Auth(email, password string) (int, error) {
	return 0, nil
}

func (m *UserLayer) CheckExists(id int) (bool, error) {
	return false, nil
}
