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
	t.Parallel()
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

func TestGistView(t *testing.T) {
	msn := newTestMission(t)
	ts := httptest.NewServer(msn.paths())
	defer ts.Close()

	tsts := []struct {
		name     string
		url      string
		wantCode int
		wantBody string
	}{
		{
			name:     "valid gist",
			url:      "/gist/1",
			wantCode: http.StatusOK,
			wantBody: "test title",
		},
		{
			name:     "invalid gist",
			url:      "/gist/2",
			wantCode: http.StatusNotFound,
			wantBody: "Gist not found",
		},
	}

	for _, tst := range tsts {
		t.Run(tst.name, func(t *testing.T) {
			//code, _, body := ts.  get(t, tst.url)
			//assert.Same(t, code, tst.wantCode)
			//if tst.wantBody != "" {
			//	assert.StringHas(t, body, tst.wantBody)
			//}
		})
	}
}
