package main

import (
	"log"
	"net/http"
)

// landing function gives byte slice as a response body
func landing(a http.ResponseWriter, b *http.Request) {
	a.Write([]byte("Here's the landing space for Gists"))
}

func main() {
	//Initialize new servemux register landing as a handler
	svrMux := http.NewServeMux()
	svrMux.HandleFunc("/", landing)

	log.Print("Listening on :9000")
	err := http.ListenAndServe(":9000", svrMux)
	log.Fatal(err)
}
