package main

import (
	"net/http"
)

func (msn *mission) paths() http.Handler {
	mx := http.NewServeMux()
	fs := http.FileServer(http.Dir("./ui/static/"))
	mx.Handle("/static/", http.StripPrefix("/static", fs))

	mx.Handle("/", http.HandlerFunc(msn.landing))
	mx.Handle("/new", http.HandlerFunc(msn.gistWrite))
	mx.Handle("/get", http.HandlerFunc(msn.gistView))
	mx.Handle("/recents", http.HandlerFunc(msn.gistRecents))

	return msn.logRq(midHeaders(mx))
}
