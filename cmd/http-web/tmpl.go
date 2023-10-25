package main

import (
	"fmt"
	"gistapp.ck89.net/internal/dblayer"
	"html/template"
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
	return tm.Format("Jan 2, 2013 at 3:04pm (SGT)")
}

var fncs = template.FuncMap{
	"fmtDate": fmtDate,
}

func newTmplCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	leaves, err := filepath.Glob("./ui/html/pages/*.tmpl")

	if err != nil {
		return nil, err
	}

	for idx, leaf := range leaves {
		fmt.Println(
			"iteration ",
			idx,
			leaf,
		)
		leafName := filepath.Base(leaf)
		fmt.Println("leaf name is ", leafName)

		tc, err := template.New(leafName).Funcs(fncs).ParseFiles("./ui/html/base.tmpl")

		if err != nil {
			return nil, err
		}

		tc, err = tc.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}
		tc, err = tc.ParseFiles(leaf)
		cache[leafName] = tc
		fmt.Println("cache leafname", cache[leafName].Name())
	}

	return cache, nil
}
