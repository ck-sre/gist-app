package main

import (
	"fmt"
	"gistapp.ck89.net/internal/dblayer"
	"gistapp.ck89.net/ui"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
)

type tmplData struct {
	PresentYr   int
	TmplGstList []dblayer.Gist
	TmplGst     dblayer.Gist
	Form        any
	Blink       string
	Validauthn  bool
	CSRFTkn     string
}

func fmtDate(tm time.Time) string {
	if tm.IsZero() {
		return ""
	}
	return tm.UTC().Format("Jan 2 15:04 UTC 2006")
}

var fncs = template.FuncMap{
	"fmtDate": fmtDate,
}

func newTmplCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	leaves, err := fs.Glob(ui.Files, "html/pages/*.tmpl")

	if err != nil {
		return nil, err
	}

	for _, leaf := range leaves {
		leafName := filepath.Base(leaf)

		ptrns := []string{"html/base.tmpl", "html/partials/*.tmpl", leaf}

		tc, err := template.New(leafName).Funcs(fncs).ParseFS(ui.Files, ptrns...)

		if err != nil {
			return nil, err
		}

		cache[leafName] = tc
		fmt.Println("cache leafname", cache[leafName].Name())
	}

	return cache, nil
}
