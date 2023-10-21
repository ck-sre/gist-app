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

type usrRegisterForm struct {
	Name            string `form:"name"`
	Email           string `form:"email"`
	Password        string `form:"password"`
	checker.Checker `form:"-"`
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

func (m *mission) usrRegister(a http.ResponseWriter, b *http.Request) {
	tmplData := m.newTmplData(b)
	tmplData.Form = usrRegisterForm{}
	m.render(a, b, http.StatusOK, "registeruser.tmpl", tmplData)
}

func (m *mission) usrRegPost(a http.ResponseWriter, b *http.Request) {
	var form usrRegisterForm

	err := m.dcdPostForm(b, &form)
	if err != nil {
		m.clErr(a, http.StatusBadRequest)
		return
	}

	form.CheckAttr(checker.NotEmpty(form.Name), "name", "This field cannot be blank")
	form.CheckAttr(checker.NotEmpty(form.Name), "email", "This field cannot be blank")
	form.CheckAttr(checker.StringMatches(form.Email, checker.EmailRegex), "email", "This field must be a valid email address")
	form.CheckAttr(checker.NotEmpty(form.Name), "password", "This field cannot be blank")
	form.CheckAttr(checker.CharMin(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.CheckPassed() {
		tmplData := m.newTmplData(b)
		tmplData.Form = form
		m.render(a, b, http.StatusUnprocessableEntity, "registeruser.tmpl", tmplData)
		return
	}

	err = m.usrs.Add(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, dblayer.ErrDuplicateEmail) {
			form.AddAttrError("email", "This email address is already in use")
			tmplData := m.newTmplData(b)
			tmplData.Form = form
			m.render(a, b, http.StatusUnprocessableEntity, "registeruser.tmpl", tmplData)
		} else {
			m.serverErr(a, b, err)
		}
		return
	}

	m.snMgr.Put(b.Context(), "blink", "Your registration was successful. Please log in.")

	http.Redirect(a, b, "/user/login", http.StatusSeeOther)

}

func (m *mission) usrLogin(a http.ResponseWriter, b *http.Request) {
	fmt.Fprintf(a, "This is a user login form")

}

func (m *mission) usrLoginPost(a http.ResponseWriter, b *http.Request) {
	fmt.Fprintf(a, "This is a user login creation")
}

func (m *mission) usrSignout(a http.ResponseWriter, b *http.Request) {
	fmt.Fprintf(a, "This is a user signout")
}
