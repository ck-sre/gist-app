package main

import (
	"context"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
)

func (msn *mission) logRq(nxt http.Handler) http.Handler {
	return http.HandlerFunc(func(a http.ResponseWriter, b *http.Request) {
		var (
			url     = b.URL.RequestURI()
			remAddr = b.RemoteAddr
			mth     = b.Method
			proto   = b.Proto
		)
		msn.logger.Info("request", "url", url, "remoteAddr", remAddr, "method", mth, "proto", proto)
		nxt.ServeHTTP(a, b)
	})
}

func midHeaders(nxt http.Handler) http.Handler {
	return http.HandlerFunc(func(a http.ResponseWriter, b *http.Request) {
		a.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; font-src fonts.gstatic.com")
		a.Header().Set("X-Frame-Options", "deny")
		a.Header().Set("X-XSS-Protection", "0")
		a.Header().Set("X-Content-Type-Options", "nosniff")
		a.Header().Set("Referrer-Policy", "origin-when-cross-origin")

		nxt.ServeHTTP(a, b)
	})
}

func (msn *mission) resurrectPanic(nxt http.Handler) http.Handler {
	return http.HandlerFunc(func(a http.ResponseWriter, b *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				a.Header().Set("Connection", "close")
				msn.serverErr(a, b, fmt.Errorf("%v", err))
			}
		}()
		nxt.ServeHTTP(a, b)
	})
}

func (msn *mission) needAuthn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(a http.ResponseWriter, b *http.Request) {
		if !msn.validAuthn(b) {
			msn.snMgr.Put(b.Context(), "pathAfterAuthn", b.URL.RequestURI())
			http.Redirect(a, b, "/usr/signin", http.StatusSeeOther)
			return
		}
		a.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(a, b)
	})
}

func noSurf(nxt http.Handler) http.Handler {
	csrfHndlr := nosurf.New(nxt)
	csrfHndlr.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
	})
	return csrfHndlr
}

func (msn *mission) authn(nxt http.Handler) http.Handler {
	return http.HandlerFunc(func(a http.ResponseWriter, b *http.Request) {
		uId := msn.snMgr.GetInt(b.Context(), "authnUserID")
		if uId == 0 {
			nxt.ServeHTTP(a, b)
			return
		}

		isPresent, err := msn.usrs.CheckExists(uId)
		if err != nil {
			msn.serverErr(a, b, err)
			return
		}

		if isPresent {
			ctx := context.WithValue(b.Context(), validCtxKey, true)
			b = b.WithContext(ctx)
		}
		nxt.ServeHTTP(a, b)

	})
}
