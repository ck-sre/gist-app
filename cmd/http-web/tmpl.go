package main

import (
	"gistapp.ck89.net/internal/dblayer"
)

//TODO: Add caching

type tmplData struct {
	TmplGst     dblayer.Gist
	TmplGstList []dblayer.Gist
	//TmplCache  map[string]*template.Template{}
}
