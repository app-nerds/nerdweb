package middlewares

import "net/http"

type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (sr *statusRecorder) Header() http.Header {
	return sr.ResponseWriter.Header()
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	return sr.ResponseWriter.Write(b)
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.Status = code
	sr.ResponseWriter.WriteHeader(code)
}
