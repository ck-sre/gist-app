package main

import (
	"errors"
	"fmt"
	"gistapp.ck89.net/internal/checker"
	"gistapp.ck89.net/internal/dblayer"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type gistWriteForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	//AttrErrors map[string]string
	checker.Checker `form:"-"` //ignore this field during decoding
}

// landing function gives byte slice as a response body
func (m *mission) landing(a http.ResponseWriter, b *http.Request) {

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

	var form gistWriteForm
	err := m.dcdPostForm(b, &form)
	if err != nil {
		m.clErr(a, http.StatusBadRequest)
		return
	}

	form.CheckAttr(checker.NotEmpty(form.Title), "title", "This field cannot be blank")
	form.CheckAttr(checker.LimitChars(form.Title, 100), "title", "This field cannot be more than 100 characters")
	form.CheckAttr(checker.NotEmpty(form.Content), "content", "This field cannot be blank")
	form.CheckAttr(checker.AllowedVal(form.Expires, 1, 365), "expires", "This field must be a number between 1 and 365")

	if !form.CheckPassed() {
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

	m.snMgr.Put(b.Context(), "blink", "Your gist has been saved successfully!")

	http.Redirect(a, b, fmt.Sprintf("/get/%d", gistid), http.StatusSeeOther)

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
