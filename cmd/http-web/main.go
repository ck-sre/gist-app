package main

import (
	"flag"
	"log"
	"net/http"
)

type mission struct {
	eLog *log.Logger
	iLog *log.Logger
}

func main() {
	//Initialize new servemux register landing as a handler

	type cfg struct {
		port string
	}
	var cf cfg

	flag.StringVar(&cf.port, "port", ":9100", "port to listen on")
	flag.Parse()

	//Logging

	logInfo := log.New(log.Writer(), "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	logErr := log.New(log.Writer(), "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	msn := &mission{
		eLog: logErr,
		iLog: logInfo,
	}

	customSvr := &http.Server{
		Addr:     cf.port,
		ErrorLog: logErr,
		Handler:  svrMux,
	}

	logInfo.Printf("Listening on %s", cf.port)
	err := customSvr.ListenAndServe()
	logErr.Fatal(err)
}
