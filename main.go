package main

import (
	"log"
	"net/http"
)

// landing function gives byte slice as a response body
func landing(a http.ResponseWriter, b *http.Request) {
	if b.URL.Path == "/" {
		http.NotFound(a, b)
		return
	}
	a.Write([]byte("Here's the landing space for Gists"))
}

func gistWrite(a http.ResponseWriter, b *http.Request) {
	a.Write([]byte("Here's a new gist we're writing"))
}

func gistView(a http.ResponseWriter, b *http.Request) {
	a.Write([]byte("This is a gist"))
}

func main() {
	//Initialize new servemux register landing as a handler
	svrMux := http.NewServeMux()
	svrMux.HandleFunc("/", landing)
	svrMux.HandleFunc("/new", gistWrite)
	svrMux.HandleFunc("/get", gistView)

	log.Print("Listening on :9000")
	err := http.ListenAndServe(":9000", svrMux)
	log.Fatal(err)
}
