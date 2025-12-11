package handlers

import (
	"log"
	"net/http"
)

type AppHandler func(http.ResponseWriter, *http.Request) error

type responseWriterWithStatus struct {
	w           http.ResponseWriter
	wroteHeader bool
}

func (r *responseWriterWithStatus) Header() http.Header {
	return r.w.Header()
}
func (r *responseWriterWithStatus) Write(b []byte) (int, error) {
	r.wroteHeader = true
	return r.w.Write(b)
}
func (r *responseWriterWithStatus) WriteHeader(statusCode int) {
	r.wroteHeader = true
	r.w.WriteHeader(statusCode)
}

func ErrorHandler(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriterWithStatus{w: w}
		err := h(rw, r)
		if err != nil {
			log.Printf("HTTP %s %s: %v", r.Method, r.URL.Path, err)
			if !rw.wroteHeader {
				http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
			}
		}
	}
}
