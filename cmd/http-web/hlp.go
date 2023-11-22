package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
	"net/http"
	"runtime/debug"
	"time"
)

func (m *mission) serverErr(w http.ResponseWriter, b *http.Request, err error) {
	var (
		method = b.Method
		url    = b.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	m.logger.Error(err.Error(), "method", method, "url", url, "stack", trace)
	//stackTrace := fmt.Sprintf("%s\n", err.Error(), debug.Stack())
	//m.eLog.Output(2, stackTrace)
	//m.logger.Error(err.Error(), "method", w.Method, "url", w.URL.RequestURI(), "stack", stackTrace)
	//m.logger.Error(stackTrace)

	if m.debug {
		http.Error(w, fmt.Sprintf("%s\n%s", err, trace), http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

// For client specific errors
func (m *mission) clErr(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For not found errors
func (m *mission) noFound(w http.ResponseWriter) {
	m.clErr(w, http.StatusNotFound)
}

func (m mission) render(w http.ResponseWriter, r *http.Request, status int, pagename string, tmplData tmplData) {

	tc, ok := m.tmplCache[pagename]
	if !ok {
		err := fmt.Errorf("This template %s  does not exist", pagename)
		m.serverErr(w, r, err)
		return
	}

	bf := new(bytes.Buffer)

	if err := tc.ExecuteTemplate(bf, "base", tmplData); err != nil {

		m.serverErr(w, r, err)
		return
	}

	w.WriteHeader(status)
	bf.WriteTo(w)
}

func (m mission) newTmplData(a *http.Request) tmplData {
	return tmplData{
		PresentYr:  time.Now().Year(),
		Blink:      m.snMgr.PopString(a.Context(), "blink"),
		Validauthn: m.validAuthn(a),
		CSRFTkn:    nosurf.Token(a),
	}
}

func (m *mission) dcdPostForm(a *http.Request, destination any) error {

	err := a.ParseForm()
	if err != nil {
		return err
	}

	err = m.formDcdr.Decode(destination, a.PostForm)

	var errDcdr *form.InvalidDecoderError
	if errors.As(err, &errDcdr) {
		panic(err)
	}

	return nil
}

func (m *mission) validAuthn(a *http.Request) bool {
	validAuthn, ok := a.Context().Value(validCtxKey).(bool)
	if !ok {
		return false
	}
	return validAuthn
}
