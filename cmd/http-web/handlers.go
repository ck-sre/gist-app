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
		"./ui/html/lowestlayer.tmpl",
		"./ui/html/partials/redirect.tmpl",
		"./ui/html/pages/landing.tmpl",
	}

	ps, err := template.ParseFiles(tmpls...)
	if err != nil {
		m.serverErr(a, err)
		return
	}
	err = ps.ExecuteTemplate(a, "lowestlayer", nil)
	if err != nil {
		m.eLog.Print(err.Error())
		m.serverErr(a, err)
	}
}

func (m *mission) gistWrite(a http.ResponseWriter, b *http.Request) {
	if b.Method != http.MethodPost {
		a.Header().Set("Allow", http.MethodPost)
		//a.Header().Set("Content-Type", "application/json")
		m.clErr(a, http.StatusMethodNotAllowed)
		return
	}

	a.Write([]byte(`{"Responseey": "Here's a new gist we're writing"}`))

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
	//a.Header().Del("Cache-Control")
	//fmt.Println(a.Header().Values("Cache-Control"))
	//a.Write([]byte(`{"ResponseKey": "This is a gist"}`))
	fmt.Fprint(a, "This is a gist with a specific id %d..", gistId)

}
