package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (m *mission) serverErr(w http.ResponseWriter, err error) {
	stackTrace := fmt.Sprintf("%s\n", err.Error(), debug.Stack())
	m.eLog.Output(2, stackTrace)
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
