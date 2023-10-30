package main

import (
	"bytes"
	"gistapp.ck89.net/internal/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMidHeaders(t *testing.T) {

	respR := httptest.NewRecorder()

	nr, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	nxt := http.HandlerFunc(func(a http.ResponseWriter, b *http.Request) {
		a.Write([]byte("pong"))
	})

	midHeaders(nxt).ServeHTTP(respR, nr)
	respRResult := respR.Result()

	//assert.Same(t, http.StatusOK, respRResult.StatusCode)

	expV := "default-src 'self'; style-src 'self'; font-src fonts.gstatic.com"
	assert.Same(t, respRResult.Header.Get("Content-Security-Policy"), expV)

	expV = "origin-when-cross-origin"
	assert.Same(t, respRResult.Header.Get("Referrer-Policy"), expV)

	expV = "nosniff"
	assert.Same(t, respRResult.Header.Get("X-Content-Type-Options"), expV)

	expV = "deny"
	assert.Same(t, respRResult.Header.Get("X-Frame-Options"), expV)

	expV = "0"
	assert.Same(t, respRResult.Header.Get("X-XSS-Protection"), expV)

	assert.Same(t, respRResult.StatusCode, http.StatusOK)

	defer respRResult.Body.Close()
	bd, err := io.ReadAll(respRResult.Body)
	if err != nil {
		t.Fatal(err)
	}

	bd = bytes.TrimSpace(bd)
	assert.Same(t, string(bd), "pong")

}
