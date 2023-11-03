package mocks

import (
	"gistapp.ck89.net/internal/dblayer"
	"time"
)

var mockGist = dblayer.Gist{
	ID:        1,
	Title:     "Mocking title",
	Content:   "Mocking content",
	CreatedOn: time.Now(),
	ExpiresOn: time.Now().AddDate(0, 0, 10),
}

type Gistdblayer struct{}

func (g *Gistdblayer) Add(title, content string, expiresOn int) (int, error) {
	return 2, nil
}

func (g *Gistdblayer) Retrieve(id int) (dblayer.Gist, error) {
	switch id {
	case 1:
		return mockGist, nil
	default:
		return dblayer.Gist{}, dblayer.ErrNoRecord
	}
}

func (g *Gistdblayer) Recent() ([]dblayer.Gist, error) {
	return []dblayer.Gist{mockGist}, nil
}
