package main

import (
	"net/http"
)

func (msn *mission) paths() {
	mx := http.NewServeMux()
	fs := http.FileServer(http.Dir("./ui/static/"))
	mx.Handle("/static/", http.StripPrefix("/static", fs))
	mx.Handle("/", http.HandlerFunc(msn.landing))
	mx.Handle("/new", http.HandlerFunc(msn.gistWrite))
	mx.Handle("/get", http.HandlerFunc(msn.gistView))
	return mx
}
