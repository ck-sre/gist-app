package main

import (
	"bytes"
	"gistapp.ck89.net/internal/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {

	respR := httptest.NewRecorder()

	nr, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ping(respR, nr)
	respRResult := respR.Result()
	assert.Same(t, http.StatusOK, respRResult.StatusCode)

	defer respRResult.Body.Close()
	bd, err := io.ReadAll(respRResult.Body)
	if err != nil {
		t.Fatal(err)
	}

	bd = bytes.TrimSpace(bd)
	assert.Same(t, string(bd), "pong")

}
