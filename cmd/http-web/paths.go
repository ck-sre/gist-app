package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (msn *mission) paths() http.Handler {

	rtr := httprouter.New()

	fs := http.FileServer(http.Dir("./ui/static/"))
	rtr.Handler(http.MethodGet, "/static/", http.StripPrefix("/static", fs))

	rtr.HandlerFunc(http.MethodGet, "/", msn.landing)
	rtr.HandlerFunc(http.MethodGet, "/new", msn.gistWrite)
	rtr.HandlerFunc(http.MethodPost, "/new", msn.gistWriteNote)
	rtr.HandlerFunc(http.MethodGet, "/get/:id", msn.gistView)

	stdMid := alice.New(msn.resurrectPanic, msn.logRq, midHeaders)
	return stdMid.Then(rtr)

	//return msn.resurrectPanic(msn.logRq(midHeaders(mx)))
}
