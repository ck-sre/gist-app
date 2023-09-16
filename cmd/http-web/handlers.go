package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// landing function gives byte slice as a response body
func (m *mission) landing(a http.ResponseWriter, b *http.Request) {
	if b.URL.Path != "/" {
		m.noFound(a)
		return
	}

	tmpls := []string{
		"./ui/html/baselayer.tmpl",
		"./ui/html/partials/redirect.tmpl",
		"./ui/html/pages/landing.tmpl",
	}

	ps, err := template.ParseFiles(tmpls...)
	if err != nil {
		m.serverErr(a, b, err)
		return
	}
	err = ps.ExecuteTemplate(a, "baselayer", nil)
	if err != nil {
		//m.eLog.Print(err.Error())
		m.serverErr(a, b, err)
		//m.logger.Error(err.Error(), "method", b.Method, "uri", b.URL.RequestURI())
		//http.Error(a, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (m *mission) gistWrite(a http.ResponseWriter, b *http.Request) {
	if b.Method != http.MethodPost {
		a.Header().Set("Allow", http.MethodPost)
		m.clErr(a, http.StatusMethodNotAllowed)
		return
	}

	title := "No suspension"
	content := "Noone wants a suspension when a new gist is created, its for the help of humanity"
	expires := 7

	gistid, err := m.gists.Add(title, content, expires)
	if err != nil {
		m.serverErr(a, b, err)
		return
	}

	http.Redirect(a, b, fmt.Sprintf("/get?gistid=%d", gistid), http.StatusSeeOther)

	a.Write([]byte(`{"Response": "Here's a new gist we're writing"}`))

}

func (m *mission) gistView(a http.ResponseWriter, b *http.Request) {

	gistId, err := strconv.Atoi(b.URL.Query().Get("gistid"))
	if err != nil || gistId < 1 {
		m.noFound(a)
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
	fmt.Fprintf(a, "This is a gist with a specific id %d..", gistId)

}
