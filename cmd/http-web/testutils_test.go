package main

import (
	"gistapp.ck89.net/internal/dblayer/mocks"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"io"
	"log/slog"
	"testing"
	"time"
)

func newTestMission(t *testing.T) *mission {
	tmplCache, err := newTmplCache()
	if err != nil {
		t.Fatal(err)
	}

	frmDcdr := form.NewDecoder()
	snMngr := scs.New()
	snMngr.Lifetime = 12 * time.Hour
	snMngr.Cookie.Secure = true

	return &mission{
		logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		gists:     &mocks.Gistdblayer{},
		usrs:      &mocks.UserLayer{},
		tmplCache: tmplCache,
		formDcdr:  frmDcdr,
		snMgr:     snMngr,
	}
}
