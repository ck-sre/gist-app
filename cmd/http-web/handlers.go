package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// landing function gives byte slice as a response body
func landing(a http.ResponseWriter, b *http.Request) {
	if b.URL.Path != "/" {
		http.NotFound(a, b)
		return
	}

	tmpls := []string{
		"./ui/html/lowestlayer.tmpl",
		"./ui/html/partials/redirect.tmpl",
		"./ui/html/pages/landing.tmpl",
	}

	ps, err := template.ParseFiles(tmpls...)
	if err != nil {
		log.Print(err.Error())
		http.Error(a, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = ps.ExecuteTemplate(a, "lowestlayer", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(a, "Internal Server Error after execution", http.StatusInternalServerError)
	}
}

func gistWrite(a http.ResponseWriter, b *http.Request) {
	if b.Method != http.MethodPost {
		a.Header().Set("Allow", http.MethodPost)
		//a.Header().Set("Content-Type", "application/json")
		http.Error(a, "That's a wrong method for this path", http.StatusMethodNotAllowed)
		return
	}

	a.Write([]byte(`{"Responseey": "Here's a new gist we're writing"}`))

}

func gistView(a http.ResponseWriter, b *http.Request) {

	gistId, err := strconv.Atoi(b.URL.Query().Get("gistid"))
	if err != nil || gistId < 1 {
		http.NotFound(a, b)
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
