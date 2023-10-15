package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (msn *mission) paths() http.Handler {

	rtr := httprouter.New()

	rtr.NotFound = http.HandlerFunc(func(a http.ResponseWriter, b *http.Request) {
		msn.noFound(a)
	})

	fs := http.FileServer(http.Dir("./ui/static/"))
	rtr.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fs))

	dyn := alice.New(msn.snMgr.LoadAndSave)

	rtr.Handler(http.MethodGet, "/", dyn.ThenFunc(msn.landing))
	rtr.Handler(http.MethodGet, "/new", dyn.ThenFunc(msn.gistWrite))
	rtr.Handler(http.MethodPost, "/new", dyn.ThenFunc(msn.gistWriteNote))
	rtr.Handler(http.MethodGet, "/get/:id", dyn.ThenFunc(msn.gistView))

	stdMid := alice.New(msn.resurrectPanic, msn.logRq, midHeaders)
	return stdMid.Then(rtr)

	//return msn.resurrectPanic(msn.logRq(midHeaders(mx)))
}
