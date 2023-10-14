package main

import (
	"errors"
	"fmt"
	"gistapp.ck89.net/internal/dblayer"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

type gistWriteForm struct {
	Title      string
	Content    string
	Expires    int
	AttrErrors map[string]string
}

// landing function gives byte slice as a response body
func (m *mission) landing(a http.ResponseWriter, b *http.Request) {
	//panic("this is a panic that's deliberate")

	tmplGsts, err := m.gists.Recent()
	if err != nil {
		m.serverErr(a, b, err)
		return
	}

	tmplData := m.newTmplData(b)
	tmplData.TmplGstList = tmplGsts

	m.render(a, b, http.StatusOK, "homebase.tmpl", tmplData)

}

func (m *mission) gistWrite(a http.ResponseWriter, b *http.Request) {
	gstData := m.newTmplData(b)
	gstData.Form = gistWriteForm{
		Expires: 365,
	}
	m.render(a, b, http.StatusOK, "write.tmpl", gstData)

}

func (m *mission) gistWriteNote(a http.ResponseWriter, b *http.Request) {

	b.Body = http.MaxBytesReader(a, b.Body, 4096)
	err := b.ParseForm()
	if err != nil {
		m.clErr(a, http.StatusBadRequest)
		return
	}

	//title := b.PostForm.Get("title")
	//content := b.PostForm.Get("content")

	expires, err := strconv.Atoi(b.PostForm.Get("expires"))
	if err != nil {
		m.clErr(a, http.StatusBadRequest)
		return
	}

	form := gistWriteForm{
		Title:      b.PostForm.Get("title"),
		Content:    b.PostForm.Get("content"),
		Expires:    expires,
		AttrErrors: make(map[string]string),
	}

	//fErr := make(map[string]string)
	if strings.TrimSpace(form.Title) == "" {
		form.AttrErrors["title"] = "This≤ field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.AttrErrors["title"] = "This field cannot be more than 100 characters"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.AttrErrors["content"] = "This field cannot be blank"
	}

	if form.Expires < 1 || expires > 365 {
		form.AttrErrors["expires"] = "This field must be a number between 1 and 365"
	}

	if len(form.AttrErrors) > 0 {
		gstData := m.newTmplData(b)
		gstData.Form = form
		m.render(a, b, http.StatusBadRequest, "write.tmpl", gstData)
		return
	}

	gistid, err := m.gists.Add(form.Title, form.Content, form.Expires)
	if err != nil {
		m.serverErr(a, b, err)
		return
	}

	http.Redirect(a, b, fmt.Sprintf("/get/%d", gistid), http.StatusSeeOther)
	//g
	//a.Write([]byte(`{"Response": "Here's a new gist we're writing"}`))

}

func (m *mission) gistView(a http.ResponseWriter, b *http.Request) {

	args := httprouter.ParamsFromContext(b.Context())
	gistId, err := strconv.Atoi(args.ByName("id"))
	if err != nil || gistId < 1 {
		m.noFound(a)
		return
	}

	gst, err := m.gists.Retrieve(gistId)
	if err != nil {
		if errors.Is(err, dblayer.ErrNoRecord) {
			m.noFound(a)
		} else {
			m.serverErr(a, b, err)
		}
		return
	}

	tmplData := m.newTmplData(b)
	tmplData.TmplGst = gst

	m.render(a, b, http.StatusOK, "viewlayer.tmpl", tmplData)

}

func (m *mission) gistRecents(a http.ResponseWriter, b *http.Request) {

	gsts, err := m.gists.Recent()
	if err != nil {
		if errors.Is(err, dblayer.ErrNoRecord) {
			m.noFound(a)
		} else {
			m.serverErr(a, b, err)
		}
		return
	}

	a.Header().Set("Content-Type", "application/json")
	a.Header().Set("Cache-Control", "public, max-age=12345600")
	a.Header().Add("Cache-Control", "public")
	a.Header().Add("Cache-Control", "max-age=12345600")
	a.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
	//a.Header().Del("Cache-Control")
	//fmt.Println(a.Header().Values("Cache-Control"))
	//a.Write([]byte(`{"ResponseKey": "This is a gist"}`))
	fmt.Fprintf(a, "+%v", gsts)

}
