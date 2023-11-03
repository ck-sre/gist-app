package dblayer

import (
	"database/sql"
	"errors"
	"time"
)

type Gist struct {
	ID        int
	Title     string
	Content   string
	CreatedOn time.Time
	ExpiresOn time.Time
}

type GistModelIface interface {
	Add(title, content string, expiresOn int) (int, error)
	Retrieve(id int) (Gist, error)
	Recent() ([]Gist, error)
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

func (g *Gistdblayer) Retrieve(id int) (Gist, error) {
	qstmt := `SELECT id, title, content, created, expires FROM gists
    WHERE expires > UTC_TIMESTAMP() AND id = ?`
	rrow := g.DB.QueryRow(qstmt, id)

	var gst Gist

	err := rrow.Scan(&gst.ID, &gst.Title, &gst.Content, &gst.CreatedOn, &gst.ExpiresOn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Gist{}, ErrNoRecord

		} else {
			return Gist{}, err
		}
	}
	return gst, nil
}

func (g *Gistdblayer) Recent() ([]Gist, error) {
	qstmt := `SELECT id, title, content, created, expires FROM gists
    WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 5`
	mrows, err := g.DB.Query(qstmt)
	if err != nil {
		return nil, err
	}

	defer mrows.Close()
	var gsts []Gist
	for mrows.Next() {
		var gst Gist
		err = mrows.Scan(&gst.ID, &gst.Title, &gst.Content, &gst.CreatedOn, &gst.ExpiresOn)
		if err != nil {
			return nil, err
		}
		gsts = append(gsts, gst)
	}
	if err = mrows.Err(); err != nil {
		return nil, err
	}
	return gsts, nil
}
