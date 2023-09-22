package main

import (
	"fmt"
	"gistapp.ck89.net/internal/dblayer"
	"html/template"
	"path/filepath"
)

type tmplData struct {
	TmplGstList []dblayer.Gist
	TmplGst     dblayer.Gist
}

func newTmplCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	leaves, err := filepath.Glob("./ui/html/pages/*.tmpl")

	if err != nil {
		return nil, err
	}

	for _, leaf := range leaves {
		fmt.Println(
			"iteration 1",
			leaf,
		)
		leafName := filepath.Base(leaf)
		fmt.Println(leafName)
		//tmplFiles := []string{
		//	"./ui/html/base.tmpl",
		//	"./ui/html/partials/redirect.tmpl",
		//	leaf,
		//}
		//
		tc, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}
		fmt.Println(tc)

		cache[leafName] = tc
		fmt.Println("cache leafname", cache[leafName])
		//leafName := filepath.Base(leaf)
		//tc, err := template.ParseFiles("./ui/html/base.tmpl")
		//
		//if err != nil {
		//	return nil, err
		//}
		//
		tc, err = tc.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}
		fmt.Println("Run tc second")
		fmt.Println(tc)
		//
		//tc, err = template.ParseFiles(leaf)
		//if err != nil {
		//	return nil, err
		//}
		//
		cache[leafName] = tc
	}

	return cache, nil
}
