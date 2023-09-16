package dblayer

import (
	"database/sql"
	"time"
)

type Gist struct {
	ID        int
	Title     string
	Content   string
	CreatedOn time.Time
	ExpiresOn time.Time
}

type Gistdblayer struct {
	DB *sql.DB
}

func (g *Gistdblayer) Add(title, content string, expiresOn int) (int, error) {
	//return 0, nil
	qstring := `INSERT INTO gists (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	qresult, err := g.DB.Exec(qstring, title, content, expiresOn)
	if err != nil {
		return 0, nil
	}

	gistid, err := qresult.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(gistid), nil
}

func (g *Gistdblayer) Retrieve(id int) (*Gist, error) {
	return nil, nil
}

func (g *Gistdblayer) Recent() ([]Gist, error) {
	return nil, nil
}
