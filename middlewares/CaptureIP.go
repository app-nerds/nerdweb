package middlewares

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type captureIP struct {
	handler http.Handler
}

func (c *captureIP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "ip", realIP(r))
	r = r.WithContext(ctx)
	c.handler.ServeHTTP(w, r)
}

/*
CaptureIP captures the caller's IP address and puts it into the
context as "ip". Example:

  mux := nerdweb.NewServeMux()
  mux.HandleFunc("/endpoint", handler)

  mux.Use(middlewares.CaptureIP())
*/
func CaptureIP() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler := &captureIP{handler: next}
			handler.ServeHTTP(w, r)
		})
	}
}
