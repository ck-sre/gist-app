package main

import (
	"gistapp.ck89.net/internal/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserRegister(t *testing.T) {
	msn := newTestMission(t)
	ts := httptest.NewServer(msn.paths())
	defer ts.Close()

	//_,_,body := ts.get(t, "/user/register")
	t.Logf("body should have been here %s", "body")
}

func TestPing(t *testing.T) {

	msn := newTestMission(t)

	ts := newTServer(t, msn.paths())
	defer ts.Close()

	stCode, _, body := ts.retrieve(t, "/ping")
	assert.Same(t, stCode, http.StatusOK)
	assert.StringHas(t, body, "pong")

}

func TestGistView(t *testing.T) {
	msn := newTestMission(t)
	ts := newTServer(t, msn.paths())
	defer ts.Close()

	tsts := []struct {
		name     string
		url      string
		wantCode int
		wantBody string
	}{
		{
			name:     "valid gist",
			url:      "/get/1",
			wantCode: http.StatusOK,
			wantBody: "test title",
		},
		{
			name:     "invalid gist",
			url:      "/get/14",
			wantCode: http.StatusNotFound,
			//wantBody: "Gist not found",
		},
		{
			name:     "Negative ID",
			url:      "/get/-1",
			wantCode: http.StatusNotFound,
			//wantBody: "Gist not found",
		},
		{
			name:     "Decimal ID",
			url:      "/get/1.23",
			wantCode: http.StatusNotFound,
			//wantBody: "Gist not found",
		},
		{
			name:     "String ID",
			url:      "/get/foo",
			wantCode: http.StatusNotFound,
			//wantBody: "Gist not found",
		},
		{
			name:     "Empty ID",
			url:      "/get/",
			wantCode: http.StatusNotFound,
			//wantBody: "Gist not found",
		},
	}

	for _, tst := range tsts {
		t.Run(tst.name, func(t *testing.T) {
			code, _, body := ts.retrieve(t, tst.url)
			assert.Same(t, code, tst.wantCode)
			if tst.wantBody != "" {
				assert.StringHas(t, body, tst.wantBody)
			}
		})
	}
}
