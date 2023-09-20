package main

import (
	"net/http"
	"runtime/debug"
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
