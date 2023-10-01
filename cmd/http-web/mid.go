package main

import "net/http"

func (msn *mission) logRq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(a http.ResponseWriter, b *http.Request) {
		var (
			url     = b.URL.RequestURI()
			remAddr = b.RemoteAddr
			mth     = b.Method
			proto   = b.Proto
		)
		msn.logger.Info("request", "url", url, "remoteAddr", remAddr, "method", mth, "proto", proto)
		next.ServeHTTP(a, b)
	})
}

func midHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(a http.ResponseWriter, b *http.Request) {
		a.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; font-src fonts.gstatic.com")
		a.Header().Set("X-Frame-Options", "deny")
		a.Header().Set("X-XSS-Protection", "0")
		a.Header().Set("X-Content-Type-Options", "nosniff")
		a.Header().Set("Referrer-Policy", "origin-when-cross-origin")

		next.ServeHTTP(a, b)
	})
}
