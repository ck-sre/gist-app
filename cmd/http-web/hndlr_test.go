package main

import (
	"gistapp.ck89.net/internal/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestUserRegister(t *testing.T) {
	msn := newTestMission(t)
	ts := newTServer(t, msn.paths())
	defer ts.Close()

	_, _, body := ts.retrieve(t, "/usr/register")
	//t.Logf("body should have been here %s", body)
	csrftTkn := getCSRFToken(t, body)

	const (
		rightName  = "Marley"
		rightPwd   = "testMarleypwd"
		rightEmail = "marley@snakeoil.com"
		rightTag   = "<form action='/usr/register' method='post' novalidate>"
	)

	tsts := []struct {
		name        string
		userName    string
		userPwd     string
		userEmail   string
		csrfToken   string
		wantCode    int
		wantBody    string
		wantFormTag string
	}{
		{
			name:      "right user",
			userName:  rightName,
			userEmail: rightEmail,
			userPwd:   rightPwd,
			csrfToken: csrftTkn,
			wantCode:  http.StatusSeeOther,
		},
		{
			name:      "invalid csrf Token",
			userName:  rightName,
			userEmail: rightEmail,
			userPwd:   rightPwd,
			csrfToken: "invalid",
			wantCode:  http.StatusBadRequest,
		},
		{
			name:      "empty name",
			userEmail: rightEmail,
			userPwd:   rightPwd,
			csrfToken: csrftTkn,
		},
	}

	for _, tst := range tsts {
		t.Run(tst.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tst.userName)
			form.Add("email", tst.userEmail)
			form.Add("password", tst.userPwd)
			form.Add("csrf", tst.csrfToken)
			code, _, body := ts.postForm(t, "/usr/register", form)
			assert.Same(t, code, tst.wantCode)
			if tst.wantBody != "" {
				assert.StringHas(t, body, tst.wantBody)
			}
			if tst.wantFormTag != "" {
				assert.StringHas(t, body, tst.wantFormTag)
			}
		})
	}

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
