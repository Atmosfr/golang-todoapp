package core_http_response

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

var StatusCodeUninitialized = -1

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     StatusCodeUninitialized,
	}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.ResponseWriter.WriteHeader(code)
	rw.statusCode = code
}

func (rw *ResponseWriter) StatusCodeOrPanic() int {
	if rw.statusCode == StatusCodeUninitialized {
		panic("status code not set")
	}
	return rw.statusCode
}
