package main

import (
	"log"
	"net/http"
)

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
