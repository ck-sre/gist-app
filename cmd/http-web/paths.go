package main

import (
	"gistapp.ck89.net/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (msn *mission) paths() http.Handler {

	rtr := httprouter.New()

	rtr.NotFound = http.HandlerFunc(func(a http.ResponseWriter, b *http.Request) {
		msn.noFound(a)
	})

	fs := http.FileServer(http.FS(ui.Files))

	rtr.Handler(http.MethodGet, "/static/*filepath", fs)

	rtr.Handler(http.MethodGet, "/ping", http.HandlerFunc(ping))

	dyn := alice.New(msn.snMgr.LoadAndSave, noSurf, msn.authn)

	rtr.Handler(http.MethodGet, "/", dyn.ThenFunc(msn.landing))
	rtr.Handler(http.MethodGet, "/info", dyn.ThenFunc(msn.info))
	rtr.Handler(http.MethodGet, "/get/:id", dyn.ThenFunc(msn.gistView))
	rtr.Handler(http.MethodGet, "/usr/register", dyn.ThenFunc(msn.usrRegister))
	rtr.Handler(http.MethodPost, "/usr/register", dyn.ThenFunc(msn.usrRegPost))
	rtr.Handler(http.MethodGet, "/usr/signin", dyn.ThenFunc(msn.usrSignin))
	rtr.Handler(http.MethodPost, "/usr/signin", dyn.ThenFunc(msn.usrSigninPost))

	protec := dyn.Append(msn.needAuthn)

	rtr.Handler(http.MethodGet, "/new", protec.ThenFunc(msn.gistWrite))
	rtr.Handler(http.MethodPost, "/new", protec.ThenFunc(msn.gistWriteNote))
	rtr.Handler(http.MethodGet, "/usr/view", protec.ThenFunc(msn.usrView))
	rtr.Handler(http.MethodGet, "/usr/signout", protec.ThenFunc(msn.usrSignoutPost))

	stdMid := alice.New(msn.resurrectPanic, msn.logRq, midHeaders)
	return stdMid.Then(rtr)

	//return msn.resurrectPanic(msn.logRq(midHeaders(mx)))
}
