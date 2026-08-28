package response

import "net/http"

var statusCodeUnitialized = -1

type Writer struct {
	http.ResponseWriter

	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *Writer {
	return &Writer{
		ResponseWriter: w,
		statusCode:     statusCodeUnitialized,
	}
}

func (rw *Writer) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *Writer) GetStatusCode() int {
	if rw.statusCode == statusCodeUnitialized {
		panic("no status code set")
	}
	return rw.statusCode
}
